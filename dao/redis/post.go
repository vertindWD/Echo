package redis

import (
	"Echo/models"
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func CreatePost(postID int64) error {
	ctx := context.Background()
	pipeline := rdb.TxPipeline()

	// 帖子发帖时间
	pipeline.ZAdd(ctx, KeyPostTimeZSet, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子的初始分数（跟发帖时间一致）
	pipeline.ZAdd(ctx, KeyPostScoreZSet, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	_, err := pipeline.Exec(ctx)
	return err
}

// GetPostIDsInOrder 根据排序规则，从 Redis 获取一页的帖子 ID
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	ctx := context.Background()

	// 默认按时间排序
	key := KeyPostTimeZSet
	if p.Order == "score" {
		key = KeyPostScoreZSet
	}

	// 确定查询的索引起点和终点 (Redis ZRange 的索引从 0 开始)
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1

	// ZRevRange 按分数从大到小查询
	return rdb.ZRangeArgs(ctx, redis.ZRangeArgs{
		Key:   key,
		Start: start,
		Stop:  end,
		Rev:   true, // 开启反向排序
	}).Result()
}

// GetPostVoteData 批量获取帖子列表的分数和当前用户的投票状态
func GetPostVoteData(ids []string, userID int64) (voteNums []int64, directions []int8, err error) {
	if len(ids) == 0 {
		return nil, nil, nil
	}

	ctx := context.Background()
	pipeline := rdb.Pipeline()

	// 提前规划好接收结果的变量
	voteNums = make([]int64, len(ids))
	directions = make([]int8, len(ids))

	// 第一遍遍历：把所有要执行的命令塞进 Pipeline
	var scoreCmds []*redis.FloatCmd
	var dirCmds []*redis.FloatCmd

	for _, id := range ids {
		// 查分数
		scoreCmds = append(scoreCmds, pipeline.ZScore(ctx, KeyPostScoreZSet, id))
		// 查当前用户的投票状态
		dirCmds = append(dirCmds, pipeline.ZScore(ctx, KeyPostVotedZSetPrefix+id, strconv.FormatInt(userID, 10)))
	}

	// 一次性把所有的命令发给 Redis 执行
	_, err = pipeline.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, nil, err
	}

	// 第二遍遍历：从 Pipeline 的执行结果中提取数据
	for i := 0; i < len(ids); i++ {
		// 提取分数
		v, _ := scoreCmds[i].Result()
		voteNums[i] = int64(v)

		// 提取投票状态 (强转为 int8)
		d, _ := dirCmds[i].Result()
		directions[i] = int8(d)
	}

	return voteNums, directions, nil
}
