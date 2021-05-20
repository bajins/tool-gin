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
	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
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
