package redis

import (
	"context"
	"fmt"

	"Echo/pkg/jwt"
)

func SetRefreshToken(userID int64, token string) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%d", KeyRefreshTokenPrefix, userID)
	return rdb.Set(ctx, key, token, jwt.RefreshTokenTTL).Err()
}

func GetRefreshToken(userID int64) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s%d", KeyRefreshTokenPrefix, userID)
	return rdb.Get(ctx, key).Result()
}

func DelRefreshToken(userID int64) error {
	ctx := context.Background()
	key := fmt.Sprintf("%s%d", KeyRefreshTokenPrefix, userID)
	return rdb.Del(ctx, key).Err()
}
