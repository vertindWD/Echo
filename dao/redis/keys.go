package redis

const (
	KeyPrefix              = "echo:"
	KeyPostTimeZSet        = "post:time"
	KeyPostScoreZSet       = "post:score"
	KeyPostVotedZSetPrefix = "post:voted:" // 记录用户及投票类型;参数是postid
)
