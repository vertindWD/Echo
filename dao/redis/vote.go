package redis

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	OneWeekInSeconds = 7 * 24 * 3600
	ScorePerVote     = 432
)

// voteScript 将读-算-写封装为一个原子操作，消除并发竞态
//
// KEYS[1] = KeyPostTimeZSet
// KEYS[2] = KeyPostScoreZSet
// KEYS[3] = KeyPostVotedZSetPrefix + postID
// ARGV[1] = postID        (ZSet 中的 member)
// ARGV[2] = userID        (voted ZSet 中的 member)
// ARGV[3] = value         (新投票方向: 1 / 0 / -1)
// ARGV[4] = 当前 Unix 时间戳
// ARGV[5] = OneWeekInSeconds
// ARGV[6] = ScorePerVote
var voteScript = redis.NewScript(`
local post_time = redis.call('ZSCORE', KEYS[1], ARGV[1])
if not post_time then
    return redis.error_reply('post not found')
end

local is_expired = (tonumber(ARGV[4]) - tonumber(post_time)) > tonumber(ARGV[5])
local ov = tonumber(redis.call('ZSCORE', KEYS[3], ARGV[2])) or 0
local value = tonumber(ARGV[3])

local dir
if value == ov then
    dir = 0
else
    dir = value
end

if dir == ov then
    return 0
end

local score_change = (dir - ov) * tonumber(ARGV[6])

if not is_expired then
    redis.call('ZINCRBY', KEYS[2], score_change, ARGV[1])
end

if dir == 0 then
    redis.call('ZREM', KEYS[3], ARGV[2])
else
    redis.call('ZADD', KEYS[3], dir, ARGV[2])
end

return 1
`)

func VoteForPost(userID, postID string, value float64) error {
	ctx := context.Background()
	keys := []string{
		KeyPostTimeZSet,
		KeyPostScoreZSet,
		KeyPostVotedZSetPrefix + postID,
	}
	err := voteScript.Run(ctx, rdb, keys,
		postID,
		userID,
		value,
		time.Now().Unix(),
		OneWeekInSeconds,
		ScorePerVote,
	).Err()

	if err == nil || err == redis.Nil {
		return nil
	}
	if strings.Contains(err.Error(), "post not found") {
		return errors.New("帖子不存在")
	}
	return err
}
