package redis

import (
	"Echo/settings"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

// Init 初始化 Redis 连接
func Init(cfg *settings.Redis) (err error) {
	// 读取参数
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis 探针回传失败: %v", err)
	}

	return nil
}

func Close() {
	if rdb != nil {
		_ = rdb.Close()
	}
}
