package logic

import (
	"Echo/dao/mysql"
	"Echo/dao/redis"
	"Echo/models"
	"Echo/pkg/kafka"
	"Echo/pkg/snowflakeID"
	"context"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	p.PostID = snowflakeID.GenID()

	// 1. 核心持久化
	if err = mysql.CreatePost(p); err != nil {
		return err
	}
	// 2. 排行榜索引
	if err = redis.CreatePost(p.PostID); err != nil {
		return err
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()

		event := map[string]interface{}{
			"type":    "create_post",
			"post_id": p.PostID,
			"user_id": p.AuthorID,
			"ts":      time.Now().Unix(),
		}
		_ = kafka.SendEvent(ctx, strconv.FormatInt(p.AuthorID, 10), event)
	}()

	return nil
}

func GetPostById(pid int64) (data *models.Post, err error) {
	return mysql.GetPostById(pid)
}

func GetPostListNew(userID int64, p *models.ParamPostList) ([]*models.PostDetail, error) {
	// 1. 去 Redis 获取一页排好序的 ID 列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder return 0 data")
		return nil, nil
	}

	// 2. 根据 ID 列表去 MySQL 查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return nil, err
	}

	// 3. 批量查询作者名
	authorIDs := make([]int64, 0, len(posts))
	for _, post := range posts {
		authorIDs = append(authorIDs, post.AuthorID)
	}
	authorMap, err := mysql.GetUsersByIDs(authorIDs)
	if err != nil {
		zap.L().Error("mysql.GetUsersByIDs failed", zap.Error(err))
		return nil, err
	}

	// 4. 去 Redis 批量查询这批帖子的点赞分数和当前用户的投票状态
	voteNums, directions, err := redis.GetPostVoteData(ids, userID)
	if err != nil {
		zap.L().Error("redis.GetPostVoteData failed", zap.Error(err))
		return nil, err
	}

	// 5. 开始拼装最终的 API 数据
	data := make([]*models.PostDetail, 0, len(posts))
	for i, post := range posts {
		postDetail := &models.PostDetail{
			Post:          post,
			AuthorName:    authorMap[post.AuthorID],
			VoteNum:       voteNums[i],
			VoteDirection: directions[i],
		}
		data = append(data, postDetail)
	}

	return data, nil
}
