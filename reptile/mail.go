/**
 *
 * @Description:
 * @Author: https://www.bajins.com
 * @File: mail.go
 * @Version: 1.0.0
 * @Time: 2019/9/16 11:36
 * @Project: tool-gin
 * @Package:
 * @Software: GoLand
 */
package reptile

import (
	"math"
	"time"
	"tool-gin/utils"
)

// 获取邮箱号后缀
func LinShiYouXiangSuffix() (string, error) {
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
	return suffixArray[s], err
}

var LIN_SHI_YOU_XIANG = "https://www.linshiyouxiang.net"

// 获取邮箱号
// prefix： 邮箱前缀
func LinShiYouXiangApply(prefix string) (map[string]interface{}, error) {
	url := LIN_SHI_YOU_XIANG + "/api/v1/mailbox/keepalive"
	param := map[string]string{
		"force_change": "1",
		"mailbox":      prefix,
		"_ts":          utils.ToString(math.Round(float64(time.Now().Unix() / 1000))),
	}
	stu, err := utils.JsonToMap(utils.HttpRequest("GET", url, param, nil))
	return stu, err
}

// 获取邮件列表
// prefix： 邮箱前缀
func LinShiYouXiangList(prefix string) string {
	url := LIN_SHI_YOU_XIANG + "/api/v1/mailbox/" + prefix
	response := utils.HttpRequest("GET", url, nil, nil)
	return response
}

// 获取邮件内容
// prefix： 邮箱前缀
// id：		邮件编号
func LinShiYouXiangGetMail(prefix, id string) string {
	url := LIN_SHI_YOU_XIANG + "/mailbox/" + prefix + "/" + id + "/source"
	response := utils.HttpRequest("GET", url, nil, nil)
	return response
}

// 删除邮件
// prefix： 邮箱前缀
// id:  	邮件编号
func LinShiYouXiangDelete(prefix, id string) string {
	url := LIN_SHI_YOU_XIANG + "/api/v1/mailbox/" + prefix + "/" + id
	response := utils.HttpRequest("DELETE", url, nil, nil)
	return response
}
