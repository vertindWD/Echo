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
	"Echo/controller"
	"Echo/dao/mysql"
	"Echo/dao/redis"
	"Echo/logger"
	"Echo/pkg/snowflakeID"
	"Echo/router"
	"Echo/settings"
	"fmt"

	"go.uber.org/zap"
)

func main() {
	// 加载配置
	if err := settings.Init(); err != nil {
		panic(err)
	}
	// 初始化日志
	if err := logger.Init(settings.Conf.Log, settings.Conf.App.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync()
	// 初始化翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init validator trans failed, err:%v\n", err)
		return
	}
	// 雪花算法生成分布式ID
	if err := snowflakeID.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	// 初始化数据库
	if err := mysql.Init(settings.Conf.MySQL); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	// 初始化redis
	if err := redis.Init(settings.Conf.Redis); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	r := router.SetupRouter()
	addr := fmt.Sprintf(":%d", settings.Conf.App.Port)
	fmt.Printf("Server is running on %s\n", addr)

	r.Run(addr)
}
