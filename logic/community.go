package logic

import (
	"Echo/dao/mysql"
	"Echo/models"
)

func GetCommunityList() ([]*models.Community, error) {
	// 查数据库
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (*models.Community, error) {
	return mysql.GetCommunityDetailByID(id)
}
