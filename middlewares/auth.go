package middlewares

import (
	"Echo/pkg/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于 JWT 的认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头的 Authorization 字段获取 Token
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "请求头中 auth 为空，请先登录"})
			c.Abort() // 物理阻断，绝不允许往下执行业务
			return
		}

		// 2. 格式校验：标准格式为 "Bearer {Token}"
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "请求头中 auth 格式有误"})
			c.Abort()
			return
		}

		// 3. 调用jwt包解析 Token
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "无效的 Token 或已过期"})
			c.Abort()
			return
		}

		// 4. 查验通过
		c.Set("userID", mc.UserID)

		c.Next()
	}
}
