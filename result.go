package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 请求返回成功
func SuccessJSON(c *gin.Context, msg string, data interface{}) {
	//var d = &struct {
	//	code    int
	//	message string
	//	data    interface{}
	//}{code:code,message:msg,data:data}
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": msg, "data": data})
	return
}

// 请求返回错误
func ErrorJSON(c *gin.Context, code int, msg string) {
	//var d = &struct {
	//	code    int
	//	message string
	//	data    interface{}
	//}{code:code,message:msg,data:""}
	c.JSON(http.StatusOK, gin.H{"code": code, "message": msg, "data": ""})
	return
}

// 默认系统错误
func SystemErrorJSON(c *gin.Context, code int, msg string) {
	//var d = &struct {
	//	code    int
	//	message string
	//	data    interface{}
	//}{code:code,message:msg,data:""}

	c.JSON(code, gin.H{"code": code, "message": msg, "data": ""})
	return
}
