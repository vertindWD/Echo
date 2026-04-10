package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	OneWeekInSeconds = 7 * 24 * 3600
	ScorePerVote     = 432
)

func VoteForPost(userID, postID string, value float64) error {
	ctx := context.Background()
	// 校验帖子时间
	postTime, err := rdb.ZScore(ctx, KeyPostTimeZSet, postID).Result()
	if err != nil {
		return errors.New("帖子不存在")
	}
	// 计算当前时间与帖子发布时间的差值
	isExpired := float64(time.Now().Unix())-postTime > OneWeekInSeconds
	ov := rdb.ZScore(ctx, KeyPostVotedZSetPrefix+postID, userID).Val()
	var dir float64
	if value == ov {
		dir = 0
	} else {
		dir = value
	}
	if dir == ov {
		return nil
	}
	scoreChange := (dir - ov) * ScorePerVote
	pipeline := rdb.TxPipeline()
	if !isExpired {
		pipeline.ZIncrBy(ctx, KeyPostScoreZSet, scoreChange, postID)
	}
	if dir == 0 {
		// 如果是取消投票，从记录中移除
		pipeline.ZRem(ctx, KeyPostVotedZSetPrefix+postID, userID)
	} else {
		// 如果是赞成(1)或反对(-1)，更新记录
		pipeline.ZAdd(ctx, KeyPostVotedZSetPrefix+postID, redis.Z{
			Score:  value,
			Member: userID,
		})
	}
	_, err = pipeline.Exec(ctx)
	return err
}
