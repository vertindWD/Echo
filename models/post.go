package models

import "time"

// Post 帖子表
type Post struct {
	// 物理主键
	ID int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	// 业务主键
	PostID      int64     `gorm:"column:post_id;uniqueIndex:idx_post_id;not null" json:"post_id,string"`                      // 帖子唯一ID（雪花算法生成）
	Title       string    `gorm:"column:title;type:varchar(128);not null" json:"title" binding:"required"`                    // 标题
	Content     string    `gorm:"column:content;type:text;not null" json:"content" binding:"required"`                        // 内容
	AuthorID    int64     `gorm:"column:author_id;index:idx_author_id;not null" json:"author_id,string"`                      // 作者ID
	CommunityID int64     `gorm:"column:community_id;index:idx_community_id;not null" json:"community_id" binding:"required"` // 所属社区ID
	Status      int32     `gorm:"column:status;type:tinyint;default:1" json:"status"`                                         // 帖子状态（1:正常 2:精华 3:置顶等）
	CreateTime  time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`                                       // 创建时间
	UpdateTime  time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`                                       // 更新时间
}

type PostDetail struct {
	*Post
	AuthorName    string `json:"author_name"`    // 作者名字
	VoteNum       int64  `json:"vote_num"`       // 帖子的热度分数
	VoteDirection int8   `json:"vote_direction"` // 当前登录用户的投票状态 (1:赞成, -1:反对, 0:未投票)
}

// TableName 指定表名
func (Post) TableName() string {
	return "post"
}
