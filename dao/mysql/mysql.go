package mysql

import (
	"fmt"
	"time"

	"Echo/models"
	"Echo/settings"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 全局对外暴露的 GORM 实例
var DB *gorm.DB

// Init 接收 settings 中的 MySQL 配置结构体进行初始化
func Init(cfg *settings.MySQL) error {
	// 1. 拼接物理连接字符串 (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)

	// 2. 建立连接
	// 这里暂时用 GORM 默认的日志配置，会把 SQL 打印到控制台
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error("MySQL 物理连接建立失败", zap.Error(err))
		return err
	}

	// 自动建表
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		zap.L().Error("自动建表失败", zap.Error(err))
		return err
	}

	err = DB.AutoMigrate(&models.Community{})
	if err != nil {
		zap.L().Error("自动建表失败", zap.Error(err))
		return err
	}

	err = DB.AutoMigrate(&models.Post{})
	if err != nil {
		zap.L().Error("Post表自动建表失败", zap.Error(err))
		return err
	}

	// 3. 提取底层原生 sql.DB 对象
	sqlDB, err := DB.DB()
	if err != nil {
		zap.L().Error("获取底层 sql.DB 失败", zap.Error(err))
		return err
	}

	// 4. 配置物理 TCP 连接池
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)                                // 最大空闲连接数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)                                // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Second) // 连接最大物理存活时间

	return nil
}
