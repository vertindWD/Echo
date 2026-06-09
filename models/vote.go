package models

import "time"

type Vote struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`
	UserID     int64     `gorm:"column:user_id;uniqueIndex:idx_user_post;not null"`
	PostID     int64     `gorm:"column:post_id;uniqueIndex:idx_user_post;not null"`
	Direction  int8      `gorm:"column:direction;not null"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime"`
}

func (Vote) TableName() string {
	return "vote"
}
