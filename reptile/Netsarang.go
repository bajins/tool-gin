/**
 *
 * @Description:
 * @Author: https://www.bajins.com
 * @File: Netsarang.go
 * @Version: 1.0.0
 * @Time: 2019/9/19 11:03
 * @Project: tool-gin
 * @Package:
 * @Software: GoLand
 */
package reptile

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
	"time"
	"tool-gin/utils"
)

func SendMail(mail, product string) error {
	if mail == "" || len(mail) == 0 {
		panic("邮箱号不能为空！")
	}
	if product == "" || len(product) == 0 {
		panic("产品不能为空！")
	}

	var url string
	if product == "Xshell" {
		url = "https://www.netsarang.com/zh/xshell-download"
	}

	if product == "Xftp" {
		url = "https://www.netsarang.com/zh/xftp-download"
	}

	if product == "Xmanager" {
		url = "https://www.netsarang.com/zh/xmanager-power-suite-download"
	}

	if product == "Xshell Plus" || product == "" {
		url = "https://www.netsarang.com/zh/xshell-plus-download"
	}

	// 定义变量，用来保存爬虫的数据
	var res string

	err := ApplyRun(clickSubmitMail(url, mail, &res))
	if err != nil {
		return err
	}
	if res == "" || len(res) == 0 || !strings.Contains(res, "感谢您提交的下载我们软件的请求") {
		return errors.New("邮箱发送失败！")
	}

	return nil
}

// 点击提交邮箱
func clickSubmitMail(url, mail string, res *string) chromedp.Tasks {
	return chromedp.Tasks{
		// 浏览器下载行为，注意设置顺序，如果不是第一个会失败
		page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorDeny),
		//network.Enable(),
		//visitWeb(url),
		//doCrawler(&res),
		//Screenshot(),
		// 跳转页面
		chromedp.Navigate(url),
		chromedp.SendKeys(`input[name="user-name"]`, strings.Split(mail, "@")[0], chromedp.BySearch),
		chromedp.SendKeys(`input[name="email"]`, mail, chromedp.BySearch),
		// 点击元素
		chromedp.Click(`input[value="开始试用"][type="submit"]`, chromedp.BySearch),
		chromedp.Sleep(10 * time.Second),
		// 查找并等待可见
		chromedp.WaitVisible(`.fusion-text h1`, chromedp.BySearch),
		// 读取HTML源码
		//chromedp.OuterHTML(`.fusion-text h1`, res, chromedp.BySearch),
		chromedp.Text(`.fusion-text h1`, res, chromedp.BySearch),
		//chromedp.TextContent(`.fusion-text h1`, res, chromedp.BySearch),
		//chromedp.Title(res),
	}
}

var netsarangInfo map[string][]interface{}

// 获取下载产品信息
func DownloadNetsarang(product string) (string, error) {
	if product == "" || len(product) == 0 {
		return "", errors.New("产品不能为空")
	}
	info := netsarangInfo[product]
	// 如果数据不为空，并且日期为今天，这么做是为了避免消耗过多的性能，每天只查询一次
	if info != nil && len(info) > 1 {
		// 判断日期是否为今天
		if utils.DateEqual(time.Now(), info[0].(time.Time)) {
			return info[1].(string), nil
		}
	}

	prefix := utils.RandomLowercaseAlphanumeric(9)
	suffix, err := LinShiYouXiangSuffix()
	if err != nil {
		return "", err
	}
	_, err = LinShiYouXiangApply(prefix)
	if err != nil {
		return "", err
	}

	mail := prefix + suffix

	log.Println("邮箱号：", mail)

	err = SendMail(mail, product)
	if err != nil {
		return "", err
	}

	time.Sleep(20 * time.Second)

	mailList := LinShiYouXiangList(prefix)

	var list []map[string]interface{}
	err = json.Unmarshal([]byte(mailList), &list)
	if err != nil {
		return "", err
	}
	listLen := len(list)
	if listLen == 0 {
		log.Println(list)
		return "", errors.New("没有邮件")
	}
	mailbox := list[listLen-1]["mailbox"].(string)
	if mailbox == "" {
		return "", errors.New("邮件前缀不存在")
	}
	mailId := list[listLen-1]["id"].(string)
	if mailId == "" {
		return "", errors.New("邮件ID不存在")
	}

	// 获取最新一封邮件
	content := LinShiYouXiangGetMail(mailbox, mailId)

	// 分割取内容
	text := strings.Split(content, "AmazonSES")
	if len(text) < 2 {
		return "", errors.New("邮件内容不正确")
	}
	// 解密，邮件协议Content-Transfer-Encoding指定了base64
	htmlText, err := base64.StdEncoding.DecodeString(text[1])
	if err != nil {
		return "", err
	}
	// 解析HTML
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htmlText))
	if err != nil {
		return "", err
	}
	tokenHtml := doc.Find(`a[target="_blank"]`)

	var attributes map[string]string

	err = ApplyRun(getDownloadUrl(tokenHtml.Text(), &attributes))
	if err != nil {
		return "", err
	}

	if attributes == nil || attributes["href"] == "" {
		return "", errors.New("没有获取到url")
	}

	// 获取最终专业版产品下载地址
	url := strings.Replace(attributes["href"], ".exe", "r.exe", -1)

	// 把产品信息存储到变量
	netsarangInfo[product] = []interface{}{time.Now(), url}

	// 在s字符串中，把old字符串替换为new字符串，n表示替换的次数，小于0表示全部替换
	return url, nil
}

// 访问带token的url获取下载地址
func getDownloadUrl(url string, attributes *map[string]string) chromedp.Tasks {
	return chromedp.Tasks{
		// 浏览器下载行为，注意设置顺序，如果不是第一个会失败
		page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorDeny),
		//network.Enable(),
		//visitWeb(url),
		//doCrawler(&res),
		//Screenshot(),
		// 跳转页面
		chromedp.Navigate(url),
		// 获取属性和值
		chromedp.Attributes(`a[target='download_frame']`, attributes, chromedp.BySearch),
	}
}

func init() {
	// 第一次调用初始化
	netsarangInfo = make(map[string][]interface{})
}
