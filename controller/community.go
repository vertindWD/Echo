package controller

import (
	"Echo/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommunityHandler 获取社区列表
// @Summary 获取所有社区板块列表
// @Description 获取所有支持发帖的社区板块信息
// @Tags 社区模块
// @Accept application/json
// @Produce application/json
// @Success 200 {object} ResponseData "成功"
// @Router /community [get]
func CommunityHandler(c *gin.Context) {
	// 查询到所有的社区
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic,GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// CommunityDetailHandler 获取社区详情
// @Summary 获取指定社区详情
// @Description 根据社区 ID 获取社区的详细介绍
// @Tags 社区模块
// @Accept application/json
// @Produce application/json
// @Param id path int true "社区 ID"
// @Success 200 {object} ResponseData "成功"
// @Router /community/{id} [get]
func CommunityDetailHandler(c *gin.Context) {
	// 获取社区id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic,GetCommunityDetail() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
