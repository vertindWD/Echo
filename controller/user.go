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
	if err := c.ShouldBindJSON(&p); err != nil {
		HandleValidatorError(c, err)
		return
	}
	pair, err := logic.SignUp(&p)
	if err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, pair)
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
	var p models.ParamLoginUsername
	if err := c.ShouldBindJSON(&p); err != nil {
		HandleValidatorError(c, err)
		return
	}
	pair, err := logic.LoginByUsername(&p)
	if err != nil {
		zap.L().Error("logic.LoginByUsername failed", zap.String("username", p.Username), zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, pair)
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
	var p models.ParamLoginEmail
	if err := c.ShouldBindJSON(&p); err != nil {
		HandleValidatorError(c, err)
		return
	}
	pair, err := logic.LoginByEmail(&p)
	if err != nil {
		zap.L().Error("logic.LoginByEmail failed", zap.String("email", p.Email), zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, pair)
}

// RefreshTokenHandler 刷新 token 接口
// @Summary 刷新 Token
// @Description 用 refresh_token 换取新的 access_token 和 refresh_token（轮转）
// @Tags 用户模块
// @Accept application/json
// @Produce application/json
// @Param object body object true "refresh_token"
// @Success 200 {object} ResponseData "成功"
// @Router /refresh_token [post]
func RefreshTokenHandler(c *gin.Context) {
	var p struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		HandleValidatorError(c, err)
		return
	}
	pair, err := logic.RefreshToken(p.RefreshToken)
	if err != nil {
		zap.L().Error("logic.RefreshToken failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, pair)
}

// LogoutHandler 登出接口
// @Summary 登出
// @Description 使当前用户的 refresh token 立即失效
// @Tags 用户模块
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Success 200 {object} ResponseData "成功"
// @Router /logout [post]
func LogoutHandler(c *gin.Context) {
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	if err := logic.Logout(userID); err != nil {
		zap.L().Error("logic.Logout failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
