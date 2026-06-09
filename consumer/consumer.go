package consumer

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"Echo/dao/mysql"
	"Echo/settings"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// Start 启动后台 Kafka 消费者
func Start(cfg *settings.Kafka) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.Address},
		Topic:    cfg.Topic,
		GroupID:  "echo-consumer-group",
		MaxBytes: 10e6,
	})

	go func() {
		for {
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				zap.L().Error("读取Kafka消息失败", zap.Error(err))
				time.Sleep(time.Second)
				continue
			}
			dispatch(msg.Value)
		}
	}()
}

type baseEvent struct {
	Type string `json:"type"`
}

func dispatch(data []byte) {
	var base baseEvent
	if err := json.Unmarshal(data, &base); err != nil {
		zap.L().Error("Kafka消息解析失败", zap.Error(err))
		return
	}
	switch base.Type {
	case "vote":
		handleVoteEvent(data)
	}
}

type voteEvent struct {
	UserID int64  `json:"user_id"`
	PostID string `json:"post_id"`
	Dir    int8   `json:"dir"`
}

func handleVoteEvent(data []byte) {
	var e voteEvent
	if err := json.Unmarshal(data, &e); err != nil {
		zap.L().Error("解析vote事件失败", zap.Error(err))
		return
	}

	postID, err := strconv.ParseInt(e.PostID, 10, 64)
	if err != nil {
		zap.L().Error("vote事件post_id格式错误", zap.String("post_id", e.PostID))
		return
	}

	if e.Dir == 0 {
		if err := mysql.DeleteVote(e.UserID, postID); err != nil {
			zap.L().Error("删除投票记录失败", zap.Error(err))
		}
		return
	}

	if err := mysql.UpsertVote(e.UserID, postID, e.Dir); err != nil {
		zap.L().Error("更新投票记录失败", zap.Error(err))
	}
}
