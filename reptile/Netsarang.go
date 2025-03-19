package reptile

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

import (
	"bytes"
	"context"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
	"tool-gin/utils"
)

const NetsarangJsonUrl = "https://update.netsarang.com/json/download/process.html"

var (
	NetsarangMap     map[string]NetsarangInfo
	NetsarangProduct = [6]string{"xshell", "xftp", "xlpd", "xshellplus", "xmanager", "powersuite"}
)

type NetsarangInfo struct {
	Time time.Time
	Url  string
}

func init() {
	// 第一次调用初始化
	NetsarangMap = make(map[string]NetsarangInfo)
}

// GetInfoUrl 获取单个产品信息
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

// NetsarangGetMail 获取可用mail
func NetsarangGetMail() (string, error) {
	prefix := utils.RandomLowercaseAlphanumeric(9)
	suffix, err := LinShiYouXiangSuffix()
	if err != nil {
		return "", err
	}
	_, err = LinShiYouXiangApply(prefix)
	if err != nil {
		return "", err
	}
	mail := prefix + "@" + suffix
	log.Println("邮箱号：", mail)
	return mail, nil
}

// NetsarangDownloadAll 获取所有链接信息
func NetsarangDownloadAll() {
	mail, err := NetsarangGetMail()
	if err != nil {
		log.Println(err)
		return
	}
	for _, app := range NetsarangProduct {
		_, err = NetsarangGetInfo(mail, app)
		if err != nil {
			log.Println(err)
		}
	}
	log.Println(NetsarangMap)
}

// NetsarangGetInfo 获取链接信息
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

	time.Sleep(10 * time.Second)
	mailList, err := LinShiYouXiangList(prefix)
	if err != nil {
		return "", err
	}
	for i := 0; i < 20; i++ {
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
	mailId := mailList[len(mailList)-1]["id"].(string)
	if mailId == "" {
		return "", errors.New("邮件ID不存在")
	}
	// 获取最新一封邮件
	msg, err := LinShiYouXiangGetMail(prefix, mailId)
	if err != nil {
		return "", err
	}
	htmlText, err := DecodeMail(msg)
	if err != nil {
		return "", err
	}
	// 解析HTML
	/*doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(htmlText)))
	  if err != nil {
	      return "", err
	  }
	  href := doc.Find(`a[target="_blank"]`).Text()*/

	exp, err := regexp.Compile(`https://www\.netsarang\.com/(.*)/downloading/\?token=(.*)<br><br>This link`)
	if err != nil {
		return "", err
	}
	href := exp.FindAllStringSubmatch(string(htmlText), -1)
	if href == nil || len(href) == 0 {
		return "", errors.New("获取token链接为空")
	}
	log.Println("token链接：", href)
	url, err := NetsarangGetUrl(href[0][1], href[0][2])
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

// NetsarangSendMail 发送邮件
func NetsarangSendMail(mail, product string) error {
	if mail == "" || len(mail) == 0 {
		return errors.New("邮箱号不能为空！")
	}
	if product == "" || len(product) == 0 {
		return errors.New("产品不能为空！")
	}

	productName := ""
	switch strings.ToLower(product) {
	case "xshell":
		productName = "xshell-download"
	case "xftp":
		productName = "xftp-download"
	case "xlpd":
		productName = "xlpd-download"
	case "xmanager":
		productName = "xmanager-download"
	case "xshellplus":
		productName = "xshell-plus-download"
	case "powersuite":
		productName = "xmanager-power-suite-download"
	}
	if productName == "" {
		return errors.New("产品不匹配")
	}
	// 请求并获取发送邮件的表单
	httpClient := utils.HttpClient{
		Method:      http.MethodGet,
		UrlText:     "https://www.netsarang.com/" + productName,
		ContentType: utils.ContentTypeMFD,
		Header:      nil,
	}
	body, err := httpClient.ReadBody()
	if err != nil {
		return err
	}
	// 解析HTML
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return err
	}
	// 找到最后一个form
	form := doc.Find(`form[novalidate="novalidate"]`).Last()
	if form.Length() < 1 {
		return errors.New("没有找到提交表单")
	}
	// 查找请求数据并构造
	inputs := form.Find(`input[type="hidden"],input[type="text"],input[type="email"]`)
	if inputs.Length() < 1 {
		return errors.New("没有找到请求数据")
	}
	// 使用make函数创建一个非nil的map，nil map不能赋值
	data := make(map[string]string)
	inputs.Each(func(i int, selection *goquery.Selection) {
		name, nbl := selection.Attr("name")
		value, vbl := selection.Attr("value")
		if nbl && vbl {
			data[name] = value
		}
	})
	if data == nil {
		return errors.New("构造请求数据失败")
	}
	data["user_name"] = mail
	data["email"] = mail
	data["productName"] = productName
	log.Println("构造数据：", data)
	// 请求发送邮件
	httpClient = utils.HttpClient{
		Method:      http.MethodPost,
		UrlText:     NetsarangJsonUrl,
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

// NetsarangGetUrl 获取下载产品信息
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
	case "ko":
		language = "1"
	case "zh":
		language = "3"
	case "ru":
		language = "8"
	case "pt":
		language = "9"
	default:
		language = "en"
	}
	params := map[string]string{
		"md":       "checkDownload",
		"token":    token,
		"language": language,
	}
	return utils.HttpReadBodyJsonMap(http.MethodPost, NetsarangJsonUrl, utils.ContentTypeMFD, params, nil)
}

// NetsarangDownloadAllDP 通过ChromeDP获取所有链接信息
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

// GetInfoUrlDP 获取单个产品信息
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

// NetsarangGetMailDP 通过ChromeDP获取可用mail
func NetsarangGetMailDP() (context.Context, context.CancelFunc, string, error) {
	var mail string
	ctx, cancel := Apply(false)
	err := chromedp.Run(ctx, GetMail24MailName(&mail))
	if err != nil {
		return nil, nil, "", err
	}
	log.Println("邮箱号：", mail)
	return ctx, cancel, mail, nil
}

// NetsarangGetInfoDP 通过ChromeDP通过ChromeDP获取链接信息
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
	for {
		if mailContent != "" {
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

// NetsarangSendMailDP 通过ChromeDP发送邮件
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
	switch strings.ToLower(product) {
	case "xshell":
		url = "https://www.netsarang.com/zh/xshell-download"
	case "xftp":
		url = "https://www.netsarang.com/zh/xftp-download"
	case "xlpd":
		url = "https://www.netsarang.com/zh/Xlpd"
	case "xmanager":
		url = "https://www.netsarang.com/zh/xmanager-download"
	case "xshellplus":
		url = "https://www.netsarang.com/zh/xshell-plus-download"
	case "powersuite":
		url = "https://www.netsarang.com/zh/xmanager-power-suite-download"
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
		//page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorDeny),
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

// NetsarangGetUrlDP 通过ChromeDP获取下载产品信息
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
		//page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorDeny),
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
