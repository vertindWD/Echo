package jwt

import (
	"errors"
	"time"

	"Echo/settings"

	"github.com/golang-jwt/jwt/v5"
)

// 声明结构体
type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenToken(userID int64, username string) (string, error) {
	// 从配置文件动态计算过期时间
	expireDuration := time.Duration(settings.Conf.Auth.JwtExpire) * time.Hour

	c := MyClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			Issuer:    "echo",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	// 从 settings 中读取秘钥，强制转换为 []byte
	return token.SignedString([]byte(settings.Conf.Auth.JwtSecret))
}

// ParseToken 验证并解析 Token
func ParseToken(tokenString string) (*MyClaims, error) {
	var mc = new(MyClaims)

	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		// 解析时同样从 settings 中读取秘钥
		return []byte(settings.Conf.Auth.JwtSecret), nil
	})

	if err != nil {
		return nil, err
	}
	if token.Valid {
		return mc, nil
	}
	return nil, errors.New("invalid token")
}
