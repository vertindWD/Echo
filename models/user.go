package models

import (
	"time"
)

// User 对应数据库的 user 表
type User struct {
	// 物理主键
	ID int64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	// 业务主键，必须打上 string 补丁防止前端 JS 精度截断
	UserID int64 `gorm:"column:user_id;uniqueIndex:idx_user_id;not null" json:"user_id,string"`
	// 用户名
	Username string `gorm:"column:username;type:varchar(64);uniqueIndex:idx_username;not null" json:"username"`
	// 密码
	Password string `gorm:"column:password;type:varchar(64);not null" json:"-"`
	// 邮箱
	Email string `gorm:"column:email;type:varchar(128)" json:"email"`
	// 时间字段
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
}

// TableName 重写 GORM 默认的表名策略
func (User) TableName() string {
	return "user"
}
