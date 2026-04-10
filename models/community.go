package models

import (
	"time"
)

// Community 映射物理表 community (板块/社区)
type Community struct {
	// 物理主键
	ID int32 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	// 业务主键
	CommunityID   uint32    `gorm:"column:community_id;uniqueIndex:idx_community_id;not null" json:"community_id"`
	CommunityName string    `gorm:"column:community_name;type:varchar(128);uniqueIndex:idx_community_name;not null" json:"community_name"`
	Introduction  string    `gorm:"column:introduction;type:varchar(256);not null" json:"introduction"`
	CreateTime    time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime    time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 强制指定表名，带上 echo 前缀
func (Community) TableName() string {
	return "community"
}
