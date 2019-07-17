package main

import (
	"github.com/gin-gonic/gin"
)

type Result struct {
	code    int
	message string
	data    interface{}
}

/**
 * 请求返回成功
 *
 * @author claer www.bajins.com
 * @date 2019/6/28 12:48
 */
func (own *Result) Success(code int, msg string, data interface{}) gin.H {
	//var d = &struct {
	//	code    int
	//	message string
	//	data    interface{}
	//}{code:code,message:msg,data:data}

	//own.code=code
	//own.message=msg
	//own.data=data

	return gin.H{"code": code, "message": msg, "data": data}
}

/**
 * 请求返回错误
 *
 * @author claer www.bajins.com
 * @date 2019/6/28 12:48
 */
func (own *Result) Error(code int, msg string) gin.H {
	//var d = &struct {
	//	code    int
	//	message string
	//	data    interface{}
	//}{code:code,message:msg,data:""}

	//own.code=code
	//own.message=msg

	return gin.H{"code": code, "message": msg, "data": ""}
}
