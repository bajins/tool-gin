/**
 *
 * @Description:
 * @Author: https://www.bajins.com
 * @File: Netsarang.go
 * @Version: 1.0.0
 * @Time: 2019/9/19 11:03
 * @Project: key-gin
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
	"key-gin/utils"
	"log"
	"strings"
	"time"
)

func SendMail(mail, product string) error {
	if mail == "" || len(mail) == 0 {
		panic("邮箱号不能为空！")
	}
	if product == "" || len(product) == 0 {
		panic("产品不能为空！")
	}

	var url string
	if product == "xshell" {
		url = "https://www.netsarang.com/zh/xshell-download"
	}

	if product == "xftp" {
		url = "https://www.netsarang.com/zh/xftp-download"
	}

	if product == "xmanager-power-suite" {
		url = "https://www.netsarang.com/zh/xmanager-power-suite-download"
	}

	if product == "xshell-plus" || product == "" {
		url = "https://www.netsarang.com/zh/xshell-plus-download"
	}

	// 定义变量，用来保存爬虫的数据
	var res string

	err := Apply(clickSubmitMail(url, mail, &res))
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

func DownloadNetsarang(product string) (string, error) {
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
	find := doc.Find(`a[target="_blank"]`)
	tokenUrl := find.Text()

	// 定义变量，用来保存爬虫的数据
	var res string

	err = Apply(getDownloadUrl(tokenUrl, &res))
	if err != nil {
		return "", err
	}

	log.Println(res)

	// 在s字符串中，把old字符串替换为new字符串，n表示替换的次数，小于0表示全部替换
	return strings.Replace(res, ".exe", "r.exe", -1), nil
}

// 访问带token的url获取下载地址
func getDownloadUrl(url string, res *string) chromedp.Tasks {
	return chromedp.Tasks{
		// 浏览器下载行为，注意设置顺序，如果不是第一个会失败
		page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorDeny),
		//network.Enable(),
		//visitWeb(url),
		//doCrawler(&res),
		//Screenshot(),
		// 跳转页面
		chromedp.Navigate(url),
		// 读取HTML源码
		chromedp.OuterHTML(`a[target='download_frame']@href`, res, chromedp.BySearch),
	}
}
