package router

import (
	"Echo/controller"
	_ "Echo/docs"
	"Echo/logger"
	"Echo/middlewares"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.Use(cors.Default())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	// 注册
	v1.POST("/signup", controller.SignUpHandler)
	// 登录
	v1.POST("/login/username", controller.LoginUsernameHandler)
	v1.POST("/login/email", controller.LoginEmailHandler)

	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)
		v1.GET("/post/list", controller.GetPostListHandler)
	}

	v1Auth := v1.Group("/")
	v1Auth.Use(middlewares.JWTAuthMiddleware())

	{
		v1Auth.POST("/post", controller.CreatePostHandler)
		v1Auth.GET("/post/:id", controller.GetPostDetailHandler)
		v1Auth.POST("/post/vote", controller.PostVoteHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
