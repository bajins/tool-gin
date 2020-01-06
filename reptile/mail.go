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
	"github.com/antchfx/htmlquery"
	"io/ioutil"
	"math"
	"math/rand"
	"time"
	"tool-gin/utils"
)

const LinShiYouXiang = "https://www.linshiyouxiang.net"

// 获取邮箱号后缀
func LinShiYouXiangSuffix() (string, error) {
	var suffixArray []string
	response, err := utils.HttpRequest("GET", LinShiYouXiang, "", nil, nil)
	if err != nil {
		return "", err
	}
	root, err := htmlquery.Parse(response.Body)
	if err != nil {
		return "", err
	}
	li := htmlquery.Find(root, "//*[@id='top']/div/div/div[2]/div/div[2]/ul/li")
	for _, row := range li {
		m := htmlquery.InnerText(row)
		suffixArray = append(suffixArray, m)
	}
	return suffixArray[rand.Intn(len(suffixArray)-1)], nil
}

// 获取邮箱号
// prefix： 邮箱前缀
func LinShiYouXiangApply(prefix string) (map[string]interface{}, error) {
	url := LinShiYouXiang + "/api/v1/mailbox/keepalive"
	param := map[string]string{
		"force_change": "1",
		"mailbox":      prefix,
		"_ts":          utils.ToString(math.Round(float64(time.Now().Unix() / 1000))),
	}
	response, err := utils.HttpRequest("GET", url, "", param, nil)
	if err != nil {
		return nil, err
	}
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	stu, err := utils.JsonToMap(string(result))
	return stu, err
}

// 获取邮件列表
// prefix： 邮箱前缀
func LinShiYouXiangList(prefix string) (string, error) {
	url := LinShiYouXiang + "/api/v1/mailbox/" + prefix
	response, err := utils.HttpRequest("GET", url, "", nil, nil)
	if err != nil {
		return "", err
	}
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// 获取邮件内容
// prefix： 邮箱前缀
// id：		邮件编号
func LinShiYouXiangGetMail(prefix, id string) (string, error) {
	url := LinShiYouXiang + "/mailbox/" + prefix + "/" + id + "/source"
	response, err := utils.HttpRequest("GET", url, "", nil, nil)
	if err != nil {
		return "", err
	}
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// 删除邮件
// prefix： 邮箱前缀
// id:  	邮件编号
func LinShiYouXiangDelete(prefix, id string) (string, error) {
	url := LinShiYouXiang + "/api/v1/mailbox/" + prefix + "/" + id
	response, err := utils.HttpRequest("DELETE", url, "", nil, nil)
	if err != nil {
		return "", err
	}
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
