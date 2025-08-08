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
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
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
}

func TestGetSvpGit(t *testing.T) {
	//fmt.Println(getSvpGit())
	//fmt.Println(getSvpDP())
	//fmt.Println(getSvpDP1())
	//fmt.Println(len(strings.Split(getSvpYseAll(), "\n")))
	fmt.Println(len(strings.Split(getSvpGitAgg(), "\n")))
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
