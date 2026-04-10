package logic

import (
	"Echo/dao/redis"
	"Echo/models"
	"strconv"
)

func PostVote(userID int64, p *models.ParamVoteData) error {
	return redis.VoteForPost(strconv.FormatInt(userID, 10), p.PostID, float64(p.Direction))
}
