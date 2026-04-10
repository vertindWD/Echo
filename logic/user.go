package logic

import (
	"Echo/dao/mysql"
	"Echo/models"
	"Echo/pkg/jwt"
	"Echo/pkg/snowflakeID"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// 将明文密码转化为加盐哈希的密文
func encryptPassword(oPassword string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(oPassword), bcrypt.DefaultCost)
	return string(hashPassword), err
}

// SignUp 统筹注册业务，并实现“注册即登录”
func SignUp(p *models.ParamSignUp) (string, error) {
	// 1. 判断用户是否存在
	err := mysql.CheckUserExist(p.Username)
	if err != nil {
		return "", err
	}

	// 2. 生成UID
	userID := snowflakeID.GenID()

	// 3. 密码哈希
	hashPassword, err := encryptPassword(p.Password)
	if err != nil {
		return "", errors.New("服务器内部错误：密码处理失败")
	}

	// 4. 构建实体模型
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: hashPassword,
	}

	if err := mysql.InsertUser(user); err != nil {
		return "", err
	}

	token, err := jwt.GenToken(userID, p.Username)
	if err != nil {
		return "", errors.New("注册成功，但自动登录失败，请手动登录")
	}

	return token, nil
}

// 用户名登录
func LoginByUsername(p *models.ParamLoginUsername) (string, error) {
	// 1. 查库
	user, err := mysql.GetUserByUsername(p.Username)
	if err != nil {
		return "", errors.New("账号不存在")
	}

	// 2. 密码校验
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p.Password))
	if err != nil {
		return "", errors.New("账号或密码错误")
	}

	// 3. 登录成功，下发 Token
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return "", errors.New("服务器内部错误:token生成失败")
	}

	return token, nil
}

// 邮箱登录
func LoginByEmail(p *models.ParamLoginEmail) (string, error) {
	// 1. 查库
	user, err := mysql.GetUserByEmail(p.Email)
	if err != nil {
		return "", errors.New("账号不存在")
	}

	// 2. 密码校验
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p.Password))
	if err != nil {
		return "", errors.New("账号或密码错误")
	}

	// 3. 登录成功，下发 Token
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return "", errors.New("服务器内部错误:token生成失败")
	}

	return token, nil
}
