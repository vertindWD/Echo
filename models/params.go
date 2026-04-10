package models

// 定义请求的参数结构体
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLoginUsername 用户名登录入参
type ParamLoginUsername struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamLoginEmail 邮箱登录入参
type ParamLoginEmail struct {
	Email    string `json:"email" binding:"required,email"` // 直接利用 Gin 的 email 正则强校验
	Password string `json:"password" binding:"required"`
}

// ParamVoteData 投票入参
type ParamVoteData struct {
	PostID    string `json:"post_id,string" binding:"required"`         // 帖子ID
	Direction int    `json:"direction" binding:"required,oneof=1 0 -1"` // 点赞(1) 点踩(-1) 取消投票(0)
}

// ParamPostList 获取帖子列表 query string 参数
type ParamPostList struct {
	Page  int64  `form:"page"`  // 页码
	Size  int64  `form:"size"`  // 每页数量
	Order string `form:"order"` // 排序依据: "time" 或 "score"
}

// ParamCreatePost 发帖请求参数
type ParamCreatePost struct {
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content" binding:"required"`
	CommunityID int64  `json:"community_id" binding:"required"`
}
