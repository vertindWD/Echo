package logic

import (
	"Echo/dao/mysql"
	"Echo/dao/redis"
	"Echo/models"
	"Echo/pkg/snowflakeID"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	// 生成post id
	p.PostID = snowflakeID.GenID()
	if err = mysql.CreatePost(p); err != nil {
		return err
	}
	return redis.CreatePost(p.PostID)
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
	// 3. 去 Redis 批量查询这批帖子的点赞分数和当前用户的投票状态
	voteNums, directions, err := redis.GetPostVoteData(ids, userID)
	if err != nil {
		zap.L().Error("redis.GetPostVoteData failed", zap.Error(err))
		return nil, err
	}

	// 4. 开始拼装最终的 API 数据
	data := make([]*models.PostDetail, 0, len(posts))
	for i, post := range posts {
		postDetail := &models.PostDetail{
			Post:          post,
			VoteNum:       voteNums[i],
			VoteDirection: directions[i],
		}
		data = append(data, postDetail)
	}

	return data, nil
}
