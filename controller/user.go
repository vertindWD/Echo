package controller

import (
	"Echo/logic"
	"Echo/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SignUpHandler 用户注册接口
// @Summary 用户注册
// @Description 注册新用户，成功后直接返回 JWT Token 实现自动登录
// @Tags 用户模块
// @Accept application/json
// @Produce application/json
// @Param object body models.ParamSignUp true "注册参数"
// @Success 200 {object} ResponseData "成功"
// @Router /signup [post]
func SignUpHandler(c *gin.Context) {
	var p models.ParamSignUp

	// 1. 参数校验
	if err := c.ShouldBindJSON(&p); err != nil {
		HandleValidatorError(c, err)
		return
	}

	// 2. 业务处理
	token, err := logic.SignUp(&p)
	if err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, token)
}

// LoginUsernameHandler 用户名登录接口
// @Summary 用户名登录
// @Description 使用用户名和密码登录
// @Tags 用户模块
// @Accept application/json
// @Produce application/json
// @Param object body models.ParamLoginUsername true "登录参数"
// @Success 200 {object} ResponseData "成功"
// @Router /login/username [post]
func LoginUsernameHandler(c *gin.Context) {
	// 参数校验
	var p models.ParamLoginUsername
	if err := c.ShouldBindJSON(&p); err != nil {
		HandleValidatorError(c, err)
		return
	}

	// 业务处理
	token, err := logic.LoginByUsername(&p)
	if err != nil {
		zap.L().Error("logic.LoginByUsername failed", zap.String("username", p.Username), zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 返回响应
	ResponseSuccess(c, token)
}

// LoginEmailHandler 邮箱登录接口
// @Summary 邮箱登录
// @Description 使用邮箱和密码登录
// @Tags 用户模块
// @Accept application/json
// @Produce application/json
// @Param object body models.ParamLoginEmail true "邮箱登录参数"
// @Success 200 {object} ResponseData "成功"
// @Router /login/email [post]
func LoginEmailHandler(c *gin.Context) {
	// 参数校验
	var p models.ParamLoginEmail
	if err := c.ShouldBindJSON(&p); err != nil {
		HandleValidatorError(c, err)
		return
	}
	// 业务处理
	token, err := logic.LoginByEmail(&p)
	if err != nil {
		zap.L().Error("logic.LoginByUsername failed", zap.String("email", p.Email), zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, token)
}
