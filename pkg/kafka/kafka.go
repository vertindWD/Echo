package kafka

import (
	"Echo/settings"
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

var writer *kafka.Writer

// Init 初始化全局 Writer
func Init(cfg *settings.Kafka) {
	writer = &kafka.Writer{
		Addr:     kafka.TCP(cfg.Address),
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
		// 允许重试
		MaxAttempts: 3,
	}
}

func SendEvent(ctx context.Context, key string, value interface{}) error {
	msgBytes, _ := json.Marshal(value)
	return writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: msgBytes,
	})
}

func Close() {
	if writer != nil {
		_ = writer.Close()
	}
}
