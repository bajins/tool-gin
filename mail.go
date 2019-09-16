/**
 *
 * @Description:
 * @Author: https://www.bajins.com
 * @File: mail.go
 * @Version: 1.0.0
 * @Time: 2019/9/16 11:36
 * @Project: key-gin
 * @Package:
 * @Software: GoLand
 */
package main

import (
	"key-gin/utils"
	"log"
	"math"
	"time"
)

// 获取邮箱号后缀
func LinShiYouXiangSuffix() string {
	suffixArray := [11]string{
		"@meantinc.com",
		"@classesmail.com",
		"@powerencry.com",
		"@groupbuff.com",
		"@figurescoin.com",
		"@navientlogin.net",
		"@programmingant.com",
		"@castlebranchlogin.com",
		"@bestsoundeffects.com",
		"@vradportal.com",
		"@a4papersize.net"}
	s, err := utils.RandomNumber(1)
	if err != nil {
		log.Fatal(err)
	}
	return suffixArray[s]
}

var LIN_SHI_YOU_XIANG = "https://www.linshiyouxiang.net"

// 获取邮箱号
// prefix： 邮箱前缀
func LinShiYouXiangGet(prefix string) map[string]interface{} {
	url := LIN_SHI_YOU_XIANG + "/api/v1/mailbox/keepalive"
	param := map[string]string{
		"force_change": string(1),
		"mailbox":      prefix,
		"_ts":          utils.ToString(math.Round(float64(time.Now().Unix() / 1000))),
	}
	stu, err := utils.JsonToMap(utils.HttpRequest("GET", url, param, nil))
	if err != nil {
		log.Fatal(err)
	}
	return stu
}

// 获取邮件列表
// prefix： 邮箱前缀
func LinShiYouXiangList(prefix string) string {
	url := LIN_SHI_YOU_XIANG + "/api/v1/mailbox/" + prefix
	response := utils.HttpClient("GET", url, nil)
	return response
}

// 删除邮件
// prefix： 邮箱前缀
// id:  	邮件编号
func LinShiYouXiangDelete(prefix, id string) string {
	url := LIN_SHI_YOU_XIANG + "/api/v1/mailbox/" + prefix + "/" + id
	response := utils.HttpClient("DELETE", url, nil)
	return response
}
