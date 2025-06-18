/**
 *
 * @Description:
 * @Author: https://www.bajins.com
 * @File: reptile_test.go
 * @Version: 1.0.0
 * @Time: 2019/9/19 11:13
 * @Project: tool-gin
 * @Package:
 * @Software: GoLand
 */
package reptile

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"github.com/parnurzeal/gorequest"
	"io"
	"log"
	"regexp"
	"strconv"
	"testing"
	"time"
)

func TestApply(t *testing.T) {
	ctx, cancel := Apply(false)
	defer cancel()
	var res string
	err := chromedp.Run(ctx, AntiDetectionHeadless(), chromedp.Tasks{
		chromedp.Sleep(20 * time.Second),
		// 跳转页面
		//chromedp.Navigate("https://intoli.com/blog/not-possible-to-block-chrome-headless/chrome-headless-test.html"),
		chromedp.Navigate("https://www.pexels.com/zh-cn/new-photos?page=1"),
		// 读取HTML源码
		chromedp.InnerHTML("html", &res, chromedp.BySearch),
	})
	t.Log(err)
	t.Log(res)
	url := "https://www.pexels.com/zh-cn/photo/3584157/"
	// 新建浏览器标签页及上下文
	ctx, cancel = chromedp.NewContext(ctx, chromedp.WithTargetID(target.ID(target.CreateTarget(url).BrowserContextID)))
	defer cancel()
	err = chromedp.Run(ctx, AntiDetectionHeadless(), chromedp.Tasks{
		// 读取HTML源码
		//chromedp.InnerHTML("html", &res, chromedp.BySearch),
	})
	t.Log(err)
	t.Log(res)
}

func TestNetsarang(t *testing.T) {
	defer func() { // 捕获panic
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
		}
	}()
	NetsarangDownloadAll()
}

func TestGetSvpGit(t *testing.T) {
	getSvpGit()
}

func TestGetSvpDP1(t *testing.T) {
	defer func() { // 捕获panic
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
		}
	}()
	//content := getSvpDP1()
	//t.Log(content)

	url := "https://vpnea.com/mfjd.html"

	request := gorequest.New()
	/*jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Failed to create cookie jar: %v", err)
	}
	request.Client.Jar = jar*/
	resp, body, errs := request.Get(url).End()
	if errs != nil {
		panic(fmt.Sprintf("一次 GET 请求失败: %v", errs))
	}
	//fmt.Println("一次 GET 响应内容: ", body)
	// 解析HTML
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(body)))
	if err != nil {
		panic(fmt.Sprintf("解析HTML失败: %v", err))
	}
	// 找到最后一个form
	form := doc.Find(`.joe_comment__respond-form`).Last()
	if form.Length() == 0 {
		panic("未找到form")
	}
	caction, _ := form.Attr("action")
	//cdt, _ := form.Attr("data-type")
	//cdc, _ := form.Attr("data-coid")
	ckb, _ := form.Find(`input[name='_']`).Last().Attr("value")
	fmt.Println(caction, ckb)

	resp, body, errs = request.Post(caction+"?time="+strconv.FormatInt(time.Now().UnixNano()/1e6, 10)).
		//SetDebug(true).
		Set("referer", url).
		Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36").
		Type("form").
		Send(`{
			"author": "csbxdh",
			"mail":   "hskdcbf@gmail.com",
			"text":   utils.RandomLower(5),
			"url": "",
			"_":   "2ec89e4c7a2e4a09e1a8cfea818253de"
		}`).
		End()
	if errs != nil {
		panic(fmt.Sprintf("提交评论失败: %v", errs))
	}
	fmt.Printf("评论提交状态: %d\n", resp.StatusCode)
	fmt.Printf("Cookies: %v\n", resp.Header.Values("Set-Cookie"))
	// 输出resp.Header
	for key, values := range resp.Header {
		fmt.Printf("%s: %s\n", key, values)
	}

	// 3. 再次 GET 请求验证评论
	resp, body, errs = request.Get(url).End() // 同首次 URL
	if errs != nil {
		panic(fmt.Sprintf("二次 GET 请求失败: %v", errs))
	}
	fmt.Printf("二次 GET 响应状态: %d\n", resp.StatusCode)
	//fmt.Println("页面内容:", resp.Body)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("读取响应体失败: %v", err))
	}
	// 解析HTML
	doc, err = goquery.NewDocumentFromReader(bytes.NewReader(bodyBytes))
	if err != nil {
		panic(fmt.Sprintf("解析HTML失败: %v", err))
	}
	// 找到最后一个pre
	pre := doc.Find(`pre`).Last()
	if pre.Length() > 0 {
		t.Log(pre.Text())
	}
}

func TestUrlRegx(t *testing.T) {
	urls := []string{
		"http://www.example.com",
		"https://example.com/path?query=123",
		"www.example.com",
		"example.com",
		"example.com/path",
		"ftp://example.com",
		"192.168.1.1", // IP address
		"localhost",
		"localhost:8080",
		"subdomain.example.co.uk",
		"example.museum",
		"http://[::1]:8080",          // IPv6
		"https://[2001:db8::1]/path", //IPv6
		"www.example-.com",           // Invalid, but test edge cases
		"-example.com",               // Invalid
		"ww-example.com",             // Invalid
		"example",                    // Invalid , but test edge cases
	}
	// 不适用于有其他文本内容参杂的情况
	//urlRegex := regexp.MustCompile(`(https?://)?([\w.-]+)(:\d+)?(/[\w./?%&=-]*)?`)
	// 更宽松，兼容性更好的正则表达式：
	urlRegex := regexp.MustCompile(`(?:(?:https?|ftp)://)?(?:(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,63}|\[(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|(?:[0-9a-fA-F]{1,4}:){1,7}:|::|localhost|\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(?::\d+)?(?:[/?#]\S*)?`)

	for _, url := range urls {
		log.Println(url, "|||||||||||", urlRegex.MatchString(url))
	}
}
