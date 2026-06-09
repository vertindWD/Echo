package logic

import (
	"Echo/dao/mysql"
	"Echo/dao/redis"
	"Echo/models"
	"Echo/pkg/jwt"
	"Echo/pkg/snowflakeID"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func encryptPassword(oPassword string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(oPassword), bcrypt.DefaultCost)
	return string(hashPassword), err
}

func SignUp(p *models.ParamSignUp) (*jwt.TokenPair, error) {
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return nil, err
	}

	userID := snowflakeID.GenID()

	hashPassword, err := encryptPassword(p.Password)
	if err != nil {
		return nil, errors.New("服务器内部错误：密码处理失败")
	}

	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: hashPassword,
	}
	if err := mysql.InsertUser(user); err != nil {
		return nil, err
	}

	return genAndStoreTokenPair(userID, p.Username)
}

func LoginByUsername(p *models.ParamLoginUsername) (*jwt.TokenPair, error) {
	user, err := mysql.GetUserByUsername(p.Username)
	if err != nil {
		return nil, errors.New("账号不存在")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p.Password)); err != nil {
		return nil, errors.New("账号或密码错误")
	}
	return genAndStoreTokenPair(user.UserID, user.Username)
}

func LoginByEmail(p *models.ParamLoginEmail) (*jwt.TokenPair, error) {
	user, err := mysql.GetUserByEmail(p.Email)
	if err != nil {
		return nil, errors.New("账号不存在")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p.Password)); err != nil {
		return nil, errors.New("账号或密码错误")
	}
	return genAndStoreTokenPair(user.UserID, user.Username)
}

// RefreshToken 用 refresh token 换取新的 token 对（含轮转）
func RefreshToken(refreshToken string) (*jwt.TokenPair, error) {
	mc, err := jwt.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("无效的 Refresh Token")
	}

	stored, err := redis.GetRefreshToken(mc.UserID)
	if err != nil || stored != refreshToken {
		return nil, errors.New("Refresh Token 已失效，请重新登录")
	}

	return genAndStoreTokenPair(mc.UserID, mc.Username)
}

// Logout 删除 Redis 中的 refresh token，使其立即失效
func Logout(userID int64) error {
	return redis.DelRefreshToken(userID)
}

func genAndStoreTokenPair(userID int64, username string) (*jwt.TokenPair, error) {
	pair, err := jwt.GenTokenPair(userID, username)
	if err != nil {
		return nil, errors.New("服务器内部错误：token 生成失败")
	}
	if err := redis.SetRefreshToken(userID, pair.RefreshToken); err != nil {
		return nil, errors.New("服务器内部错误：token 存储失败")
	}
	return pair, nil
}
