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
