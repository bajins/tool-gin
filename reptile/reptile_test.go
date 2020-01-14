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
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"testing"
)

func TestApply(t *testing.T) {
	//ctx:=context.Background()
}

func TestCDP(t *testing.T) {
	// 定义变量，用来保存爬虫的数据
	var res string
	//ctx, cancel := ApplyDebug()
	//defer cancel()
	ctx, _ := Apply()
	err := chromedp.Run(ctx, Crawler(&res))
	if err != nil {
		panic(err)
	}
	t.Log(res)
	if res == "" || len(res) == 0 {
		t.Log("邮箱发送失败！")
	}
}

func Crawler(res *string) chromedp.Tasks {
	return chromedp.Tasks{
		network.Enable(),
		//visitWeb(url),
		//doCrawler(&res),
		//Screenshot(),
		// 跳转页面
		chromedp.Navigate("https://www.netsarang.com/zh/thank-you-download/"),
		// 读取HTML源码
		//chromedp.OuterHTML(`.fusion-text h1::text`, res, chromedp.BySearch),
		chromedp.Title(res),
	}
}

func TestLinShiYouXiangSuffix(t *testing.T) {
	LinShiYouXiangSuffix()
}

func TestLinShiYouXiangList(t *testing.T) {
	list, _ := LinShiYouXiangList("5wij52emu")
	t.Log(list)
}

func TestLinShiYouXiangGetMail(t *testing.T) {
	_, err := DownloadNetsarang("xshell")
	t.Log(err)

}

func TestSendMail(t *testing.T) {
	SendMail("", "xshell")
}

func TestDownloadNetsarang(t *testing.T) {

	url, err := DownloadNetsarang("xshell")
	t.Log(url)
	t.Error(err == nil)
}

func TestGetMail24(t *testing.T) {
	//GetMail24()
	var test string
	//ctx, cancel := ApplyDebug()
	//defer cancel()
	ctx, _ := ApplyDebug()
	err := chromedp.Run(ctx, GetMail24(Mail24, &test))
	t.Log(err)
	t.Log(test)
	err = chromedp.Run(ctx, GetMail24List(Mail24, &test))
	t.Log(err)
	fmt.Println(test)
}
