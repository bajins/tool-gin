package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义业务状态码常量 net/http/status.go
const (
	SUCCESS = 0
	ERROR   = http.StatusInternalServerError

	ErrInvalidParams = 1002
)

// MsgFlags 存储状态码对应的消息
var MsgFlags = map[int]string{
	SUCCESS: "success",
	ERROR:   "fail",

	ErrInvalidParams: "请求参数错误",
}

// GetMsg 根据状态码获取消息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	// 如果未找到，返回通用的服务器错误消息
	return MsgFlags[ERROR]
}

// Response 是返回给前端的统一结构体
type Response struct {
	Code    int         `json:"code"`    // 业务状态码
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 数据
}

// Result 是一个辅助函数，用于发送统一格式的响应
func Result(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: msg,
		Data:    data,
	})
	/*var d = &struct {
		code    int
		message string
		data    interface{}
	}{code:code,message:msg,data:""}*/
	//c.JSON(code, gin.H{"code": code, "message": msg, "data": ""})
}

// respond 是一个内部方法，用于发送 JSON 响应
func (ctx *Context) respond(httpStatus int, code int, msg string, data interface{}) {
	ctx.C.JSON(httpStatus, Response{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

func (ctx *Context) SuccessJSON(msg string, data interface{}) {
	ctx.respond(http.StatusOK, http.StatusOK, msg, data)
}

func (ctx *Context) ErrorJSON(code int, msg string) {
	ctx.respond(http.StatusOK, code, msg, nil)
}

func (ctx *Context) SystemErrorJSON(code int, msg string) {
	ctx.respond(http.StatusInternalServerError, code, msg, nil)
}

// ErrorByCode 方法，根据预定义的错误码发送失败响应
func (ctx *Context) ErrorByCode(code int) {
	msg := GetMsg(code)
	ctx.respond(http.StatusOK, code, msg, nil)
}

// BindAndValidate 方法，封装了参数绑定和验证逻辑
// 如果出错，它会自动发送错误响应
func (ctx *Context) BindAndValidate(obj interface{}) bool {
	if err := ctx.C.ShouldBind(obj); err != nil {
		// 如果绑定或验证失败，直接返回参数错误
		ctx.ErrorByCode(ErrInvalidParams)
		return false
	}
	return true
}
