package controller

import (
	"Echo/logic"
	"Echo/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// PostVoteHandler 帖子投票接口
// @Summary 帖子投票
// @Description 登录用户对指定的帖子进行赞成(1)、反对(-1)或取消投票(0)操作
// @Tags 帖子模块
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object body models.ParamVoteData true "投票参数"
// @Success 200 {object} ResponseData "成功"
// @Router /post/vote [post]
func PostVoteHandler(c *gin.Context) {
	// 参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(&p); err != nil {
		HandleValidatorError(c, err)
		return
	}
	uid, err := GetCurrentUserID(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 或 CodeNeedLogin
		return
	}
	if err := logic.PostVote(uid, p); err != nil {
		zap.L().Error("logic.PostVote(uid, p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
