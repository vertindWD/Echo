package main

// @title Echo 社区 API 文档
// @version 1.0
// @description 这是基于 Gin 框架开发的社区后端 API 服务。
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8080
// @BasePath /api/v1
import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Echo/consumer"
	"Echo/controller"
	"Echo/dao/mysql"
	"Echo/dao/redis"
	"Echo/logger"
	"Echo/pkg/kafka"
	"Echo/pkg/snowflakeID"
	"Echo/router"
	"Echo/settings"

	"go.uber.org/zap"
)

func main() {
	if err := settings.Init(); err != nil {
		panic(err)
	}
	if err := logger.Init(settings.Conf.Log, settings.Conf.App.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync()
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init validator trans failed, err:%v\n", err)
		return
	}
	if err := snowflakeID.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	if err := mysql.Init(settings.Conf.MySQL); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	if err := redis.Init(settings.Conf.Redis); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	kafka.Init(settings.Conf.Kafka)
	consumer.Start(settings.Conf.Kafka)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.App.Port),
		Handler: router.SetupRouter(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("服务启动失败", zap.Error(err))
		}
	}()
	zap.L().Info("服务已启动", zap.String("addr", srv.Addr))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.L().Info("正在关闭服务...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error("HTTP 服务关闭异常", zap.Error(err))
	}
	kafka.Close()
	redis.Close()

	zap.L().Info("服务已安全退出")
}
