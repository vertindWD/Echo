package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "userID"

var ErrorUserNotLogin = errors.New("用户未登录")

// GetCurrentUserID 获取当前登录的用户ID
func GetCurrentUserID(c *gin.Context) (userID int64, err error) {
	// 1. 从上下文中取出 userID
	uid, ok := c.Get(CtxUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}

	// 2. 类型断言：确保取出来的值是 int64 类型
	userID, ok = uid.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}

	return
}
