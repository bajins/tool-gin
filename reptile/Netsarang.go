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
	"context"
	"encoding/base64"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
	"tool-gin/utils"
)

var NetsarangMap map[string]NetsarangInfo

type NetsarangInfo struct {
	Time time.Time
	Url  string
}

func init() {
	// 第一次调用初始化
	NetsarangMap = make(map[string]NetsarangInfo)
}

// 获取单个产品信息
func GetInfoUrl(product string) (string, error) {
	info := NetsarangMap[product]
	if NetsarangMap == nil || info.Url == "" || len(info.Url) == 0 || !utils.DateEqual(time.Now(), info.Time) {
		mail, err := NetsarangGetMail()
		if err != nil {
			return "", err
		}
		url, err := NetsarangGetInfo(mail, product)
		if err != nil {
			return "", err
		}
		info = NetsarangInfo{Time: time.Now(), Url: url}
	}
	return info.Url, nil
}

// 获取可用mail
func NetsarangGetMail() (string, error) {
	prefix := utils.RandomLowercaseAlphanumeric(9)
	suffix, err := LinShiYouXiangSuffix()
	if err != nil {
		return "", err
	}
	res, err := LinShiYouXiangApply(prefix)
	if err != nil {
		return "", err
	}
	log.Println(res)
	mail := prefix + suffix
	log.Println("邮箱号：", mail)
	return mail, nil
}

// 获取所有链接信息
func NetsarangDownloadAll() {
	mail, err := NetsarangGetMail()
	if err != nil {
		log.Println(err)
		return
	}
	_, err = NetsarangGetInfo(mail, "Xshell")
	if err != nil {
		log.Println(err)
	}
	_, err = NetsarangGetInfo(mail, "Xftp")
	if err != nil {
		log.Println(err)
	}
	_, err = NetsarangGetInfo(mail, "Xmanager")
	if err != nil {
		log.Println(err)
	}
	_, err = NetsarangGetInfo(mail, "Xshell Plus")
	if err != nil {
		log.Println(err)
	}
	log.Println(NetsarangMap)
}

// 获取链接信息
func NetsarangGetInfo(mail, product string) (string, error) {
	if mail == "" || len(mail) == 0 {
		return "", errors.New("mail不能为空")
	}
	if product == "" || len(product) == 0 {
		return "", errors.New("product不能为空")
	}
	info := NetsarangMap[product]
	// 如果数据不为空，并且日期为今天，这么做是为了避免消耗过多的性能，每天只查询一次
	if info.Url != "" && len(info.Url) > 1 && utils.DateEqual(time.Now(), info.Time) {
		return "", nil
	}
	err := NetsarangSendMail(mail, product)
	if err != nil {
		return "", err
	}
	prefix := strings.Split(mail, "@")[0]
	mailList, err := LinShiYouXiangList(prefix)
	if err != nil {
		return "", err
	}
	for i := 0; i < 30; i++ {
		if len(mailList) > 0 {
			break
		}
		time.Sleep(10 * time.Second)
		mailList, err = LinShiYouXiangList(prefix)
		if err != nil {
			return "", err
		}
	}
	if len(mailList) == 0 {
		return "", errors.New("没有邮件")
	}
	mailbox := mailList[len(mailList)-1]["mailbox"].(string)
	if mailbox == "" {
		return "", errors.New("邮件前缀不存在")
	}
	mailId := mailList[len(mailList)-1]["id"].(string)
	if mailId == "" {
		return "", errors.New("邮件ID不存在")
	}
	// 获取最新一封邮件
	content, err := LinShiYouXiangGetMail(mailbox, mailId)
	if err != nil {
		return "", err
	}
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
	href := doc.Find(`a[target="_blank"]`).Text()

	exp, err := regexp.Compile(`https://www\.netsarang\.com/(.*)/downloading/\?token=(.*)`)
	if err != nil {
		return "", err
	}
	hrf := exp.FindAllStringSubmatch(href, -1)
	log.Println("token链接：", hrf[0][0])
	if hrf == nil || len(hrf) == 0 {
		return "", errors.New("获取token链接为空")
	}
	url, err := NetsarangGetUrl(hrf[0][1], hrf[0][2])
	if err != nil {
		return "", err
	}
	if url == nil || url["downlink"] == "" {
		return "", errors.New("没有获取到url")
	}
	// 获取最终专业版产品下载地址
	// 在s字符串中，把old字符串替换为new字符串，n表示替换的次数，小于0表示全部替换
	ur := strings.Replace(url["downlink"].(string), ".exe", "r.exe", -1)
	// 把产品信息存储到变量
	NetsarangMap[product] = NetsarangInfo{Time: time.Now(), Url: ur}
	return ur, nil
}

// 发送邮件
func NetsarangSendMail(mail, product string) error {
	if mail == "" || len(mail) == 0 {
		return errors.New("邮箱号不能为空！")
	}
	if product == "" || len(product) == 0 {
		return errors.New("产品不能为空！")
	}

	productCode := ""
	productName := ""
	if strings.EqualFold(product, "Xshell") {
		productCode = "4203"
		productName = "xshell-download"
	} else if strings.EqualFold(product, "Xftp") {
		productCode = "4242"
		productName = "xftp-download"
	} else if strings.EqualFold(product, "Xmanager") {
		productCode = "4066"
		productName = "xmanager-power-suite-download"
	} else if strings.EqualFold(product, "Xshell Plus") || product == "" {
		productCode = "4132"
		productName = "xshell-plus-download"
	}
	if productCode == "" || productName == "" {
		return errors.New("产品不匹配")
	}
	data := map[string]string{
		"_wpcf7":                "3016",
		"_wpcf7_version":        "5.1.1",
		"_wpcf7_locale":         "en_US",
		"_wpcf7_unit_tag":       "wpcf7-f3016-p" + productCode + "-o2",
		"_wpcf7_container_post": productCode,
		"g-recaptcha-response":  "",
		"md":                    "setDownload",
		"language":              "3",
		"downloadType":          "0",
		"licenseType":           "0",
		"action":                "/json/download/process.html",
		"user-name":             mail,
		"email":                 mail,
		"company":               "",
		"productName":           productName,
	}
	httpClient := utils.HttpClient{
		Method:      http.MethodPost,
		UrlText:     "https://www.netsarang.com/json/download/process.html",
		ContentType: utils.ContentTypeMFD,
		Params:      data,
		Header:      nil,
	}
	js, err := httpClient.HttpReadBodyJsonMap()
	if err != nil {
		return err
	}
	if js == nil || !js["result"].(bool) || js["errorCounter"].(float64) != 0 {
		return errors.New("邮箱发送失败！")
	}
	return nil
}

// 获取下载产品信息
func NetsarangGetUrl(lang, token string) (map[string]interface{}, error) {
	if lang == "" || len(lang) == 0 {
		return nil, errors.New("lang不能为空")
	}
	if token == "" || len(token) == 0 {
		return nil, errors.New("token不能为空")
	}
	var language string
	switch lang {
	case "en":
		language = "2"
		break
	case "ko":
		language = "1"
		break
	case "zh":
		language = "3"
		break
	case "ru":
		language = "8"
		break
	case "pt":
		language = "9"
		break
	default:
		language = "en"
		break
	}
	params := map[string]string{
		"md":       "checkDownload",
		"token":    token,
		"language": language,
	}
	return utils.HttpReadBodyJsonMap(http.MethodPost, "https://www.netsarang.com/json/download/process.html", utils.ContentTypeMFD, params, nil)
}

// 通过ChromeDP获取所有链接信息
func NetsarangDownloadAllDP() {
	ctx, cancel, mail, err := NetsarangGetMailDP()
	defer cancel()
	if err != nil {
		log.Println(err)
		return
	}
	_, err = NetsarangGetInfoDP(ctx, mail, "Xshell")
	if err != nil {
		log.Println(err)
	}
	_, err = NetsarangGetInfoDP(ctx, mail, "Xftp")
	if err != nil {
		log.Println(err)
	}
	_, err = NetsarangGetInfoDP(ctx, mail, "Xmanager")
	if err != nil {
		log.Println(err)
	}
	_, err = NetsarangGetInfoDP(ctx, mail, "Xshell Plus")
	if err != nil {
		log.Println(err)
	}
	log.Println(NetsarangMap)
}

// 获取单个产品信息
func GetInfoUrlDP(product string) (string, error) {
	info := NetsarangMap[product]
	if NetsarangMap == nil || info.Url == "" || len(info.Url) == 0 || !utils.DateEqual(time.Now(), info.Time) {
		ctx, cancel, mail, err := NetsarangGetMailDP()
		defer cancel()
		if err != nil {
			return "", err
		}
		url, err := NetsarangGetInfoDP(ctx, mail, product)
		if err != nil {
			return "", err
		}
		info = NetsarangInfo{Time: time.Now(), Url: url}
	}
	return info.Url, nil
}

// 通过ChromeDP获取可用mail
func NetsarangGetMailDP() (context.Context, context.CancelFunc, string, error) {
	var mail string
	ctx, cancel := ApplyRun()
	err := chromedp.Run(ctx, GetMail24MailName(&mail))
	if err != nil {
		return nil, nil, "", err
	}
	log.Println("邮箱号：", mail)
	return ctx, cancel, mail, nil
}

// 通过ChromeDP通过ChromeDP获取链接信息
func NetsarangGetInfoDP(ctx context.Context, mail, product string) (string, error) {
	if ctx == nil {
		return "", errors.New("context不能为空")
	}
	if mail == "" || len(mail) == 0 {
		return "", errors.New("mail不能为空")
	}
	if product == "" || len(product) == 0 {
		return "", errors.New("product不能为空")
	}
	info := NetsarangMap[product]
	// 如果数据不为空，并且日期为今天，这么做是为了避免消耗过多的性能，每天只查询一次
	if info.Url != "" && len(info.Url) > 1 && utils.DateEqual(time.Now(), info.Time) {
		return "", nil
	}
	err := NetsarangSendMailDP(ctx, mail, product)
	if err != nil {
		return "", err
	}
	var mailContent string
	err = chromedp.Run(ctx, GetMail24LatestMail(&mailContent))
	if err != nil {
		return "", err
	}
	for i := 0; i < 30; {
		if mailContent != "" || err != nil {
			break
		}
		err = chromedp.Run(ctx, GetMail24LatestMail(&mailContent))
		if err != nil {
			return "", err
		}
	}
	exp, err := regexp.Compile(`https://www\.netsarang\.com/(.*)/downloading/\?token=(.*)`)
	if err != nil {
		return "", err
	}
	hrf := exp.FindString(mailContent)
	log.Println("token链接：", hrf)
	if hrf == "" || len(hrf) == 0 {
		return "", errors.New("获取token链接为空")
	}
	url, err := NetsarangGetUrlDP(ctx, hrf)
	if err != nil {
		return "", err
	}

	// 把产品信息存储到变量
	NetsarangMap[product] = NetsarangInfo{Time: time.Now(), Url: url}
	return url, nil
}

// 通过ChromeDP发送邮件
func NetsarangSendMailDP(ctx context.Context, mail, product string) error {
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

// 通过ChromeDP获取下载产品信息
func NetsarangGetUrlDP(ctx context.Context, url string) (string, error) {
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
