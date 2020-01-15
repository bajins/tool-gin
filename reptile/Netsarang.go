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
	"context"
	"errors"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"log"
	"regexp"
	"strings"
	"time"
	"tool-gin/utils"
)

var NetsarangInfo map[string][]interface{}

func init() {
	// 第一次调用初始化
	NetsarangInfo = make(map[string][]interface{})
}

// 获取可用mail
func NetsarangGetMail() (context.Context, context.CancelFunc, string, error) {
	var mail string
	ctx, cancel := ApplyDebug()
	err := chromedp.Run(ctx, GetMail24MailName(&mail))
	if err != nil {
		return nil, nil, "", err
	}
	log.Println("邮箱号：", mail)
	return ctx, cancel, mail, nil
}

func tickerGetInfo() {
	ticker := time.NewTicker(time.Minute * 20)
	<-ticker.C
	NetsarangDownloadAll()
}

// 获取所有链接信息
func NetsarangDownloadAll() {
	ctx, cancel, mail, err := NetsarangGetMail()
	defer cancel()
	if err != nil {
		log.Println(err)
		return
	}
	_, err = NetsarangGetInfo(ctx, mail, "Xshell")
	if err != nil {
		log.Println(err)
		go tickerGetInfo()
	}
	_, err = NetsarangGetInfo(ctx, mail, "Xftp")
	if err != nil {
		log.Println(err)
		go tickerGetInfo()
	}
	_, err = NetsarangGetInfo(ctx, mail, "Xmanager")
	if err != nil {
		log.Println(err)
		go tickerGetInfo()
	}
	_, err = NetsarangGetInfo(ctx, mail, "Xshell Plus")
	if err != nil {
		log.Println(err)
		go tickerGetInfo()
	}
	log.Println(NetsarangInfo)
}

// 获取链接信息
func NetsarangGetInfo(ctx context.Context, mail, product string) (string, error) {
	if ctx == nil {
		return "", errors.New("context不能为空")
	}
	if mail == "" || len(mail) == 0 {
		return "", errors.New("mail不能为空")
	}
	if product == "" || len(product) == 0 {
		return "", errors.New("product不能为空")
	}
	info := NetsarangInfo[product]
	// 如果数据不为空，并且日期为今天，这么做是为了避免消耗过多的性能，每天只查询一次
	if info != nil && len(info) > 1 {
		// 判断日期是否为今天
		if utils.DateEqual(time.Now(), info[0].(time.Time)) {
			return "", nil
		}
	}
	err := NetsarangSendMail(ctx, mail, product)
	if err != nil {
		return "", err
	}
	var mailList string
	err = chromedp.Run(ctx, GetMail24LatestMail(&mailList))
	if err != nil {
		return "", err
	}
	for i := 0; i < 30; {
		if mailList != "" || err != nil {
			break
		}
		err = chromedp.Run(ctx, GetMail24LatestMail(&mailList))
		if err != nil {
			return "", err
		}
	}
	exp, err := regexp.Compile(`https://www\.netsarang\.com/.*/downloading/\?token=.*`)
	if err != nil {
		return "", err
	}
	hrf := exp.FindString(mailList)
	log.Println("token链接：", hrf)
	if hrf == "" {
		return "", errors.New("获取token链接为空")
	}
	url, err := NetsarangGetUrl(ctx, hrf)
	if err != nil {
		return "", err
	}
	// 把产品信息存储到变量
	NetsarangInfo[product] = []interface{}{time.Now(), url}
	return url, nil
}

// 发送邮件
func NetsarangSendMail(ctx context.Context, mail, product string) error {
	if ctx == nil {
		return errors.New("context不能为空")
	}
	if mail == "" || len(mail) == 0 {
		return errors.New("邮箱号不能为空！")
	}
	if product == "" || len(product) == 0 {
		return errors.New("产品不能为空！")
	}

	var url string
	if strings.EqualFold(product, "Xshell") {
		url = "https://www.netsarang.com/zh/xshell-download"
	}

	if strings.EqualFold(product, "Xftp") {
		url = "https://www.netsarang.com/zh/xftp-download"
	}

	if strings.EqualFold(product, "Xmanager") {
		url = "https://www.netsarang.com/zh/xmanager-power-suite-download"
	}

	if strings.EqualFold(product, "Xshell Plus") || product == "" {
		url = "https://www.netsarang.com/zh/xshell-plus-download"
	}
	if url == "" {
		return errors.New("产品不匹配，url为空")
	}
	// 定义变量，用来保存爬虫的数据
	var res string
	tags, _ := chromedp.Targets(ctx)
	if tags != nil {
		log.Println("当前浏览器实例标签页数：", len(tags))
	}
	// 新建浏览器标签页及上下文
	ctx, cancel := chromedp.NewContext(ctx, chromedp.WithTargetID(target.ID(target.CreateTarget(url).BrowserContextID)))
	defer cancel()
	err := chromedp.Run(ctx, chromedp.Tasks{
		// 浏览器下载行为，注意设置顺序，如果不是第一个会失败
		page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorDeny),
		// 跳转页面
		chromedp.Navigate(url),
		chromedp.SendKeys(`input[name="user-name"]`, strings.Split(mail, "@")[0], chromedp.BySearch),
		chromedp.SendKeys(`input[name="email"]`, mail, chromedp.BySearch),
		// 点击元素
		chromedp.Click(`input[value="开始试用"][type="submit"]`, chromedp.BySearch),
		chromedp.Sleep(20 * time.Second),
		// 查找并等待可见
		chromedp.WaitVisible(`//*[@id="content"]/div/div/div[2]/div/div/div/div[1]/h1`, chromedp.BySearch),
		chromedp.Text(`//*[@id="content"]/div/div/div[2]/div/div/div/div[1]/h1`, &res, chromedp.BySearch),
	})
	if err != nil {
		return err
	}
	if res == "" || len(res) == 0 || !strings.Contains(res, "感谢您提交的下载我们软件的请求") {
		return errors.New("邮箱发送失败！")
	}
	return nil
}

// 获取下载产品信息
func NetsarangGetUrl(ctx context.Context, url string) (string, error) {
	if ctx == nil {
		return "", errors.New("context不能为空")
	}
	if url == "" || len(url) == 0 {
		return "", errors.New("url不能为空")
	}
	var attributes map[string]string
	// 新建浏览器标签页及上下文
	ctx, cancel := chromedp.NewContext(ctx, chromedp.WithTargetID(target.ID(target.CreateTarget(url).BrowserContextID)))
	defer cancel()
	err := chromedp.Run(ctx, chromedp.Tasks{
		// 浏览器下载行为，注意设置顺序，如果不是第一个会失败
		page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorDeny),
		// 跳转页面
		chromedp.Navigate(url),
		// 获取属性和值
		chromedp.Attributes(`a[target='download_frame']`, &attributes, chromedp.BySearch),
	})
	if err != nil {
		return "", err
	}
	if attributes == nil || attributes["href"] == "" {
		return "", errors.New("没有获取到url")
	}
	// 获取最终专业版产品下载地址
	// 在s字符串中，把old字符串替换为new字符串，n表示替换的次数，小于0表示全部替换
	return strings.Replace(attributes["href"], ".exe", "r.exe", -1), nil
}
