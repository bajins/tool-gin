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
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"math"
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
	return suffixArray[utils.RandIntn(len(suffixArray)-1)], nil
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

const Mail24 = "http://24mail.chacuo.net"

func GetMail24(url string, res *string) chromedp.Tasks {
	return chromedp.Tasks{
		// 浏览器下载行为，注意设置顺序，如果不是第一个会失败
		page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorDeny),
		//network.Enable(),
		//visitWeb(url),
		//doCrawler(&res),
		//Screenshot(),
		// 跳转页面
		chromedp.Navigate(url),
		chromedp.Value("converts", res, chromedp.ByID),
		// 点击元素
		//chromedp.Click(`input[value="开始试用"][type="submit"]`, chromedp.BySearch),
		//chromedp.Sleep(20 * time.Second),
		// 查找并等待可见
		//chromedp.WaitVisible(`//*[@id="content"]/div/div/div[2]/div/div/div/div[1]/h1`, chromedp.BySearch),
		// 读取HTML源码
		//chromedp.OuterHTML(`.fusion-text h1`, res, chromedp.BySearch),
		//chromedp.Text(`//*[@id="content"]/div/div/div[2]/div/div/div/div[1]/h1`, res, chromedp.BySearch),
		//chromedp.TextContent(`.fusion-text h1`, res, chromedp.BySearch),
		//chromedp.Title(res),
	}
}
func GetMail24List(url string, res *string) chromedp.Tasks {
	return chromedp.Tasks{
		// 浏览器下载行为，注意设置顺序，如果不是第一个会失败
		page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorDeny),
		//chromedp.Sleep(20 * time.Second),
		// 读取HTML源码
		chromedp.OuterHTML(`//*[@id="mailtooltipss"]/ul/ins`, res, chromedp.BySearch),
	}
}

const GuerrillaMail = "https://www.guerrillamail.com/zh"

const TempMail = "https://temp-mail.org/zh"

const Moakt = "https://www.moakt.com/zh"

const Mail5 = "http://www.5-mail.com"

const YopMail = "http://www.yopmail.com/zh"

const MinuteMail10 = "https://10minutemail.com/10MinuteMail/index.html"

const IncognitoMail = "http://www.incognitomail.com"

const MailCatch = "http://mailcatch.com/en/disposable-email"

const MinteMail = "https://www.mintemail.com"

const Maildu = "http://www.maildu.de"

const MailDrop = "https://maildrop.cc"

const EM9 = "https://9em.org"

const CS = "https://www.cs.email/zh"
