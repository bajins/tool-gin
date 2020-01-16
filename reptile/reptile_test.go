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
	"tool-gin/utils"
)

func TestLinShiYouXiangSuffix(t *testing.T) {
	LinShiYouXiangSuffix()
}

func TestLinShiYouXiangList(t *testing.T) {
	list, _ := LinShiYouXiangList("5wij52emu")
	t.Log(list)
}

func TestGetMail24(t *testing.T) {
	//GetMail24()
	var test string
	//ctx, cancel := ApplyDebug()
	//defer cancel()
	ctx, _ := ApplyDebug()
	err := chromedp.Run(ctx, GetMail24MailName(&test))
	t.Log(err)
	t.Log(test)
	err = chromedp.Run(ctx, GetMail24LatestMail(&test))
	t.Log(err)
	fmt.Println(test)
}

func TestApply(t *testing.T) {
	url := "https://www.netsarang.com/zh/downloading/?token=d1hXUC05Y3RVWXhJNWt6NF9rUHhDQUBaVjZXVkJRQU51VHEtRi1PVm1MQUFR"
	response, err := utils.HttpReadBodyString("GET", url, "", nil, nil)
	fmt.Println(response, err)
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
