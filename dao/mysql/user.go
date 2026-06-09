package mysql

import (
	"Echo/models"
	"errors"
)

func InsertUser(u *models.User) error {
	err := DB.Create(u).Error
	return err
}

func CheckUserExist(username string) error {
	var count int64
	err := DB.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户已存在")
	}
	return nil
}

// 用户名查询
func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

// 邮箱查询
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

// GetUsersByIDs 批量查询，返回 userID -> username 的映射
func GetUsersByIDs(ids []int64) (map[int64]string, error) {
	var users []models.User
	err := DB.Select("user_id, username").Where("user_id IN (?)", ids).Find(&users).Error
	if err != nil {
		return nil, err
	}
	result := make(map[int64]string, len(users))
	for _, u := range users {
		result[u.UserID] = u.Username
	}
	return result, nil
}
