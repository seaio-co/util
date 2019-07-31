package common

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// SuccessBody 成功响应结构体
type ResponseBody struct {
	Status   string `json:"status"`
	ErrorMsg string `json:"error"`
	Data     interface{}
}

// Success 成功
func Success(data interface{}) ResponseBody {
	return ResponseBody{
		Status: "200",
		Data:   data,
	}
}

// Error 失败
func Error(errcode, msg string) ResponseBody {
	return ResponseBody{
		Status:   errcode,
		ErrorMsg: msg,
	}
}

// ReturnSuccessResponse HTTP接口返回数据信息
func ReturnSuccessResponse(ctx *gin.Context, data interface{}) {

	response := &ResponseBody{
		Status: "200",
		Data:   data,
	}

	ctx.JSON(200, response)
	return
}

// ReturnErrorResponse HTTP接口返回错误 //err返回时表示系统出现错误，errMsg返回表示请求无效
func ReturnErrorResponse(ctx *gin.Context, statusStr, errorMsg string) {

	response := &ResponseBody{
		Status:   statusStr,
		ErrorMsg: errorMsg,
	}

	status, _ := strconv.Atoi(statusStr)
	ctx.JSON(status, response)
	return
}
