package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 请求返回成功
func Success(msg string, data interface{}) gin.H {
	//var d = &struct {
	//	code    int
	//	message string
	//	data    interface{}
	//}{code:code,message:msg,data:data}

	return gin.H{"code": http.StatusOK, "message": msg, "data": data}
}

// 请求返回错误
func Error(code int, msg string) gin.H {
	//var d = &struct {
	//	code    int
	//	message string
	//	data    interface{}
	//}{code:code,message:msg,data:""}

	return gin.H{"code": code, "message": msg, "data": ""}
}

// 默认系统错误
func SystemError() gin.H {
	//var d = &struct {
	//	code    int
	//	message string
	//	data    interface{}
	//}{code:code,message:msg,data:""}

	return gin.H{"code": http.StatusInternalServerError, "message": "系统错误！", "data": ""}
}
