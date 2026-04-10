package controller

import (
	"Echo/logic"
	"Echo/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler 创建新帖子
// @Summary 创建新帖子
// @Description 登录用户可以在指定的社区板块内发布新帖子
// @Tags 帖子模块
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object body models.ParamCreatePost true "发帖的参数"
// @Success 200 {object} ResponseData "成功"
// @Router /post [post]
func CreatePostHandler(c *gin.Context) {
	// 获取参数及参数校验
	var p models.ParamCreatePost
	if err := c.ShouldBindJSON(&p); err != nil {
		HandleValidatorError(c, err)
		return
	}
	// 获取发请求的用户id
	userID, err := GetCurrentUserID(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID(c) failed, token parse error or not set")
		ResponseError(c, CodeServerBusy)
		return
	}
	post := &models.Post{
		Title:       p.Title,
		Content:     p.Content,
		CommunityID: p.CommunityID,
		AuthorID:    userID,
	}
	// 创建帖子
	if err := logic.CreatePost(post); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情
// @Summary 获取帖子详情
// @Description 根据帖子 ID 获取帖子的详细内容
// @Tags 帖子模块
// @Accept application/json
// @Produce application/json
// @Param id path int true "帖子 ID"
// @Success 200 {object} ResponseData "成功"
// @Router /post/{id} [get]
func GetPostDetailHandler(c *gin.Context) {
	// 获取参数并校验
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 根据id取出帖子数据
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取按时间或热度排序的帖子列表
// @Summary 获取帖子列表
// @Description 支持按时间和分数排序。如果有 Token，会额外返回当前用户的点赞状态。
// @Tags 帖子模块
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌 (游客选填)"
// @Param page query int false "页码"
// @Param size query int false "每页数量"
// @Param order query string false "排序规则 (time 或 score)"
// @Success 200 {object} ResponseData "成功"
// @Router /post/list [get]
func GetPostListHandler(c *gin.Context) {
	// 获取参数并校验
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: "time", // 默认按时间
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler with invalid params", zap.Error(err))
		HandleValidatorError(c, err)
		return
	}
	// 获取数据
	userID, err := GetCurrentUserID(c)
	if err != nil {
		userID = 0 // 未登录状态，下面查出来的 VoteDirection 自然全都是 0
	}

	data, err := logic.GetPostListNew(userID, p)
	if err != nil {
		zap.L().Error("logic.GetPostListNew(userID, p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}
