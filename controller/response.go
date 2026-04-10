package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}

// HandleValidatorError 统一处理 Gin 的参数绑定与校验错误
func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": CodeInvalidParam,
			"msg":  "请求参数格式错误",
			"err":  err.Error(),
		})
		return
	}
	errData := removeTopStruct(errs.Translate(trans))
	c.JSON(http.StatusBadRequest, gin.H{
		"code": CodeInvalidParam,
		"msg":  "参数校验失败",
		"Data": errData,
	})
}
