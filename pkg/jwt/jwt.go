package jwt

import (
	"errors"
	"time"

	"Echo/settings"

	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
	RefreshTokenTTL  = 7 * 24 * time.Hour
)

type MyClaims struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// GenTokenPair 生成 access token + refresh token 对
func GenTokenPair(userID int64, username string) (*TokenPair, error) {
	accessToken, err := genToken(userID, username, tokenTypeAccess,
		time.Duration(settings.Conf.Auth.JwtExpire)*time.Hour)
	if err != nil {
		return nil, err
	}
	refreshToken, err := genToken(userID, username, tokenTypeRefresh, RefreshTokenTTL)
	if err != nil {
		return nil, err
	}
	return &TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func genToken(userID int64, username, tokenType string, ttl time.Duration) (string, error) {
	c := MyClaims{
		UserID:    userID,
		Username:  username,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			Issuer:    "echo",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(settings.Conf.Auth.JwtSecret))
}

// ParseToken 解析并校验 access token，拒绝 refresh token
func ParseToken(tokenString string) (*MyClaims, error) {
	mc, err := parseRawToken(tokenString)
	if err != nil {
		return nil, err
	}
	if mc.TokenType != tokenTypeAccess {
		return nil, errors.New("invalid token type")
	}
	return mc, nil
}

// ParseRefreshToken 解析并校验 refresh token，拒绝 access token
func ParseRefreshToken(tokenString string) (*MyClaims, error) {
	mc, err := parseRawToken(tokenString)
	if err != nil {
		return nil, err
	}
	if mc.TokenType != tokenTypeRefresh {
		return nil, errors.New("invalid token type")
	}
	return mc, nil
}

func parseRawToken(tokenString string) (*MyClaims, error) {
	mc := new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		return []byte(settings.Conf.Auth.JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return mc, nil
}
