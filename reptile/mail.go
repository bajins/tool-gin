package reptile

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

import (
	"encoding/base64"
	"errors"
	"math"
	"net/http"
	"net/mail"
	"strings"
	"time"
	"tool-gin/utils"

	"github.com/antchfx/htmlquery"
	"github.com/chromedp/chromedp"
)

// DecodeMail 解码邮件内容 https://github.com/alexcesaro/quotedprintable
func DecodeMail(msg *mail.Message) ([]byte, error) {
	body := utils.BytesToStringByBuffer(msg.Body)
	if len(body) == 0 || body == "" {
		return nil, errors.New("邮件内容不正确")
	}
	encoding := msg.Header.Get("Content-Transfer-Encoding")
	// 解码，邮件协议Content-Transfer-Encoding指定了编码方式
	if encoding == "base64" {
		body, err := base64.StdEncoding.DecodeString(body)
		return body, err
	}
	return nil, errors.New("解码方式错误：" + encoding)
}

const LinShiYouXiang = "https://www.linshiyouxiang.net"

// LinShiYouXiangSuffix 获取邮箱号后缀
func LinShiYouXiangSuffix() (string, error) {
	var suffixArray []string
	response, err := utils.HttpRequest(http.MethodGet, LinShiYouXiang, "", nil, nil)
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
	suffixArrayLen := len(suffixArray)
	if suffixArrayLen == 0 {
		return "", nil
	}
	return suffixArray[utils.RandIntn(len(suffixArray)-1)], nil
}

// LinShiYouXiangApply 获取邮箱号
// prefix： 邮箱前缀
func LinShiYouXiangApply(prefix string) (map[string]interface{}, error) {
	url := LinShiYouXiang + "/api/v1/mailbox/keepalive"
	param := map[string]string{
		"force_change": "1",
		"mailbox":      prefix,
		"_ts":          utils.ToString(math.Round(float64(time.Now().Unix() / 1000))),
	}
	r, e := utils.HttpReadBodyJsonMap(http.MethodGet, url, "", param, nil)
	return r, e
}

// LinShiYouXiangList 获取邮件列表
// prefix： 邮箱前缀
func LinShiYouXiangList(prefix string) ([]map[string]interface{}, error) {
	url := LinShiYouXiang + "/api/v1/mailbox/" + prefix
	return utils.HttpReadBodyJsonMapArray(http.MethodGet, url, "", nil, nil)
}

// LinShiYouXiangGetMail 获取邮件内容
// prefix： 邮箱前缀
// id：		邮件编号
//
// 获取到邮件需要做以下操作
// 分割取内容
// text := strings.Split(content, "AmazonSES")
// 解密，邮件协议Content-Transfer-Encoding指定了base64
// htmlText, err := base64.StdEncoding.DecodeString(text[1])
// 解析HTML
// doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlText))
func LinShiYouXiangGetMail(prefix, id string) (*mail.Message, error) {
	url := LinShiYouXiang + "/mailbox/" + prefix + "/" + id + "/source"
	content, err := utils.HttpReadBodyString(http.MethodGet, url, "", nil, nil)
	if err != nil {
		return nil, err
	}
	r := strings.NewReader(content)
	m, err := mail.ReadMessage(r) // 解析邮件
	return m, err
}

// LinShiYouXiangDelete 删除邮件
// prefix： 邮箱前缀
// id:  	邮件编号
func LinShiYouXiangDelete(prefix, id string) (map[string]interface{}, error) {
	url := LinShiYouXiang + "/api/v1/mailbox/" + prefix + "/" + id
	return utils.HttpReadBodyJsonMap(http.MethodDelete, url, "", nil, nil)
}

const Mail24 = "http://24mail.chacuo.net"

func GetMail24MailName(res *string) chromedp.Tasks {
	return chromedp.Tasks{
		// 浏览器下载行为，注意设置顺序，如果不是第一个会失败
		//page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorDeny),
		//network.Enable(),
		//visitWeb(url),
		//doCrawler(&res),
		//Screenshot(),
		// 跳转页面
		chromedp.Navigate(Mail24),
		chromedp.Sleep(20 * time.Second),
		// 查找并等待可见
		chromedp.WaitVisible("mail_cur_name", chromedp.ByID),
		chromedp.WaitReady("mail_cur_name", chromedp.ByID),
		chromedp.Value("mail_cur_name", res, chromedp.ByID),
		// 点击元素
		//chromedp.Click(`input[value="开始试用"][type="submit"]`, chromedp.BySearch),
		// 读取HTML源码
		//chromedp.OuterHTML(`.fusion-text h1`, res, chromedp.BySearch),
		//chromedp.Text(`//*[@id="content"]/div/div/div[2]/div/div/div/div[1]/h1`, res, chromedp.BySearch),
		//chromedp.TextContent(`.fusion-text h1`, res, chromedp.BySearch),
		//chromedp.Title(res),
	}
}

// GetMail24List 获取邮件列表
func GetMail24List(res *string) chromedp.Tasks {
	return chromedp.Tasks{
		// 浏览器下载行为，注意设置顺序，如果不是第一个会失败
		//page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorDeny),
		chromedp.Sleep(20 * time.Second),
		// 读取HTML源码
		chromedp.InnerHTML(`//*[@id="convertd"]`, res, chromedp.BySearch),
	}
}

// GetMail24LatestMail 获取最新邮件
func GetMail24LatestMail(res *string) chromedp.Tasks {
	return chromedp.Tasks{
		// 浏览器下载行为，注意设置顺序，如果不是第一个会失败
		//page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorDeny),
		chromedp.WaitVisible(`//*[@id="convertd"]/tr[1]`, chromedp.BySearch),
		chromedp.Click(`//*[@id="convertd"]/tr[1]`, chromedp.BySearch),
		chromedp.Sleep(10 * time.Second),
		//chromedp.WaitVisible(`//*[@id="mailview_data"]`, chromedp.BySearch),
		chromedp.Click(`//*[@id="mailview"]/thead/tr[1]/td/a[1]`, chromedp.BySearch),
		chromedp.TextContent(`//*[@id="mailview_data"]`, res, chromedp.BySearch),
	}
}
