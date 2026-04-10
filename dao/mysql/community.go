package mysql

import "Echo/models"

func GetCommunityList() ([]*models.Community, error) {
	var communityList []*models.Community
	err := DB.Find(&communityList).Error
	return communityList, err
}

func GetCommunityDetailByID(id int64) (community *models.Community, err error) {
	community = new(models.Community)
	// 使用 GORM 查询第一条匹配的数据
	err = DB.Where("community_id = ?", id).First(community).Error
	return community, err
}
