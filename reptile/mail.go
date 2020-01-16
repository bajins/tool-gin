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
	"math"
	"net/http"
	"time"
	"tool-gin/utils"
)

const LinShiYouXiang = "https://www.linshiyouxiang.net"

// 获取邮箱号后缀
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
	return utils.HttpReadBodyJsonMap(http.MethodGet, url, "", param, nil)
}

// 获取邮件列表
// prefix： 邮箱前缀
func LinShiYouXiangList(prefix string) ([]map[string]interface{}, error) {
	url := LinShiYouXiang + "/api/v1/mailbox/" + prefix
	return utils.HttpReadBodyJsonArray(http.MethodGet, url, "", nil, nil)
}

// 获取邮件内容
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
func LinShiYouXiangGetMail(prefix, id string) (string, error) {
	url := LinShiYouXiang + "/mailbox/" + prefix + "/" + id + "/source"
	return utils.HttpReadBodyString(http.MethodGet, url, "", nil, nil)
}

// 删除邮件
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
		page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorDeny),
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

// 获取邮件列表
func GetMail24List(res *string) chromedp.Tasks {
	return chromedp.Tasks{
		// 浏览器下载行为，注意设置顺序，如果不是第一个会失败
		page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorDeny),
		chromedp.Sleep(20 * time.Second),
		// 读取HTML源码
		chromedp.InnerHTML(`//*[@id="convertd"]`, res, chromedp.BySearch),
	}
}

// 获取最新邮件
func GetMail24LatestMail(res *string) chromedp.Tasks {
	return chromedp.Tasks{
		// 浏览器下载行为，注意设置顺序，如果不是第一个会失败
		page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorDeny),
		chromedp.WaitVisible(`//*[@id="convertd"]/tr[1]`, chromedp.BySearch),
		chromedp.Click(`//*[@id="convertd"]/tr[1]`, chromedp.BySearch),
		chromedp.Sleep(10 * time.Second),
		//chromedp.WaitVisible(`//*[@id="mailview_data"]`, chromedp.BySearch),
		chromedp.Click(`//*[@id="mailview"]/thead/tr[1]/td/a[1]`, chromedp.BySearch),
		chromedp.TextContent(`//*[@id="mailview_data"]`, res, chromedp.BySearch),
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
