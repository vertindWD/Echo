package logic

import (
	"Echo/dao/redis"
	"Echo/models"
	"Echo/pkg/kafka"
	"context"
	"strconv"
	"time"
)

func PostVote(userID int64, p *models.ParamVoteData) error {

	err := redis.VoteForPost(strconv.FormatInt(userID, 10), p.PostID, float64(p.Direction))

	// 利用 Kafka 的顺序性保证用户投票状态的流式处理
	if err == nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			kafka.SendEvent(ctx, p.PostID, map[string]interface{}{
				"type":    "vote",
				"user_id": userID,
				"post_id": p.PostID,
				"dir":     p.Direction,
			})
		}()
	}
	return err
}
