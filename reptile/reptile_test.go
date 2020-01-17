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
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
	"testing"
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
	//ctx, cancel := Apply(true)
	//defer cancel()
	ctx, _ := Apply(true)
	err := chromedp.Run(ctx, GetMail24MailName(&test))
	t.Log(err)
	t.Log(test)
	err = chromedp.Run(ctx, GetMail24LatestMail(&test))
	t.Log(err)
	fmt.Println(test)
}

func TestApply(t *testing.T) {
	ctx, _ := Apply(true)
	var res string
	err := chromedp.Run(ctx, chromedp.Tasks{
		AntiDetectionHeadless(),
		// 跳转页面
		//chromedp.Navigate("https://intoli.com/blog/not-possible-to-block-chrome-headless/chrome-headless-test.html"),
		chromedp.Navigate("https://www.pexels.com/zh-cn/new-photos?page=1"),
		// 读取HTML源码
		chromedp.InnerHTML("html", &res, chromedp.BySearch),
	})
	t.Log(err)
	t.Log(res)
	// 新建浏览器标签页及上下文
	ctx, cancel := chromedp.NewContext(ctx, chromedp.WithTargetID(target.ID(target.CreateTarget("https://www.pexels.com/zh-cn/photo/3584157/").BrowserContextID)))
	defer cancel()
	err = chromedp.Run(ctx, chromedp.Tasks{
		AntiDetectionHeadless(),
		// 读取HTML源码
		//chromedp.InnerHTML("html", &res, chromedp.BySearch),
	})
	t.Log(err)
	t.Log(res)
}
