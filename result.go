package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// SuccessJSON 请求返回成功
func SuccessJSON(c *gin.Context, msg string, data interface{}) {
	//var d = &struct {
	//	code    int
	//	message string
	//	data    interface{}
	//}{code:code,message:msg,data:data}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": msg, "data": data})
}

// ErrorJSON 请求返回错误
func ErrorJSON(c *gin.Context, code int, msg string) {
	//var d = &struct {
	//	code    int
	//	message string
	//	data    interface{}
	//}{code:code,message:msg,data:""}
	c.JSON(http.StatusOK, gin.H{"code": code, "message": msg, "data": ""})
}

// SystemErrorJSON 默认系统错误
func SystemErrorJSON(c *gin.Context, code int, msg string) {
	//var d = &struct {
	//	code    int
	//	message string
	//	data    interface{}
	//}{code:code,message:msg,data:""}

	c.JSON(code, gin.H{"code": code, "message": msg, "data": ""})
}
