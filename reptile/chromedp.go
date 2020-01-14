/**
 *
 * @Description:
 * @Author: https://www.bajins.com
 * @File: chromedp.go
 * @Version: 1.0.0
 * @Time: 2019/9/19 9:31
 * @Project: tool-gin
 * @Package:
 * @Software: GoLand
 */
package reptile

import (
	"context"
	"errors"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"strings"
	"time"
	"tool-gin/utils"
)

// 显示浏览器窗口启动，结束时不关闭浏览器实例
//
// context.Context部分不能抽离，否则会报 context canceled
func ApplyDebug(actions chromedp.Action) error {
	// 创建缓存目录
	//dir, err := ioutil.TempDir("", "chromedp-example")
	//if err != nil {
	//	panic(err)
	//}
	//defer os.RemoveAll(dir)

	opts := []chromedp.ExecAllocatorOption{
		// 设置UA，防止有些页面识别headless模式
		chromedp.UserAgent(utils.UserAgent),
		// 窗口最大化
		chromedp.Flag("start-maximized", true),
	}
	ctx, _ := chromedp.NewExecAllocator(context.Background(), opts...)

	// 自定义记录器
	ctx, _ = chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))

	// 设置超时时间
	ctx, _ = context.WithTimeout(ctx, 3*time.Minute)

	ctx, _ = context.WithCancel(ctx)

	// listen network event
	//listenForNetworkEvent(ctx)
	return chromedp.Run(ctx, actions)
}

// 不显示浏览器窗口启动，结束时关闭浏览器实例
//
// context.Context部分不能抽离，否则会报 context canceled
func ApplyRun(actions chromedp.Action) error {
	// 创建缓存目录
	//dir, err := ioutil.TempDir("", "chromedp-example")
	//if err != nil {
	//	panic(err)
	//}
	//defer os.RemoveAll(dir)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		// 禁用GPU，不显示GUI
		chromedp.DisableGPU,
		// 隐身模式启动
		chromedp.Flag("incognito", true),
		// 取消沙盒模式
		chromedp.NoSandbox,
		// 忽略证书错误
		chromedp.Flag("ignore-certificate-errors", true),
		// 指定浏览器分辨率
		//chromedp.WindowSize(1600, 900),
		// 窗口最大化
		chromedp.Flag("start-maximized", true),
		//
		chromedp.Flag("in-process-plugins", true),
		// 不加载图片, 提升速度
		chromedp.Flag("disable-images", true),
		// 禁用扩展
		chromedp.Flag("disable-extensions", true),
		// 隐藏滚动条, 应对一些特殊页面
		chromedp.Flag("hide-scrollbars", true),
		// 设置UA，防止有些页面识别headless模式
		chromedp.UserAgent(utils.UserAgent),
		// 设置用户数据目录
		//chromedp.UserDataDir(dir),
		//chromedp.ExecPath("C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe"),
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	// 关闭chrome实例
	defer cancel()

	// 自定义记录器
	ctx, cancel = chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
	// 释放所有资源，并等待释放结束
	defer cancel()

	// 设置超时时间
	ctx, cancel = context.WithTimeout(ctx, 3*time.Minute)
	// 超时关闭chrome实例
	defer cancel()

	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	// listen network event
	//listenForNetworkEvent(ctx)
	return chromedp.Run(ctx, actions)
}

//监听
func listenForNetworkEvent(ctx context.Context) {
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *network.EventResponseReceived:
			resp := ev.Response
			if len(resp.Headers) != 0 {
				// log.Printf("received headers: %s", resp.Headers)
				if strings.Index(resp.URL, ".ts") != -1 {
					log.Printf("received headers: %s", resp.URL)
				}
			}
		case *network.WebSocketResponse:
			respH := ev.Headers
			log.Println("WebSocketResponse", respH)
		case *network.Cookie:
			domain := ev.Domain
			log.Println("Cookie", domain)
		case *network.Headers:
			log.Println("Headers", ev)
		}
		// other needed network Event
	})
}

// 任务 主要用来设置cookie ，获取登录账号后的页面
func visitWeb(url string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctxt context.Context) error {
			expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
			// 设置cookie
			success, err := network.SetCookie("ASP.NET_SessionId", "这里是值").
				WithExpires(&expr).
				// 访问网站主体
				WithDomain(url).
				WithHTTPOnly(true).
				Do(ctxt)
			if err != nil {
				return err
			}
			if !success {
				return errors.New("无法设置cookie")
			}

			return nil
		}),
		// 页面跳转
		chromedp.Navigate(url),
	}
}

// 截图
func Screenshot() chromedp.Tasks {
	var buf []byte

	task := chromedp.Tasks{
		chromedp.CaptureScreenshot(&buf),
		chromedp.ActionFunc(func(context.Context) error {
			return ioutil.WriteFile("testimonials.png", buf, 0644)
		}),
	}
	//if err := ioutil.WriteFile("fullScreenshot.png", buf, 0644); err != nil {
	//	log.Fatal("生成图片错误：", err)
	//}
	return task
}

// 任务 主要执行翻页功能和或者html
func DoCrawler(url string, res *string) chromedp.Tasks {
	return chromedp.Tasks{
		// 浏览器下载行为，注意设置顺序，如果不是第一个会失败
		page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorDeny),
		network.Enable(),
		//visitWeb(url),
		//doCrawler(&res),
		//Screenshot(),
		// 跳转页面
		chromedp.Navigate(url),
		// 查找并等待可见
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		// 等待1秒
		chromedp.Sleep(1 * time.Second),
		// 点击元素
		//chromedp.Click(`.pagination li:nth-last-child(4) a`, chromedp.BySearch),
		// 读取HTML源码
		chromedp.OuterHTML(`body`, res, chromedp.ByQuery),
		//chromedp.Text(`.fusion-text h1`, res, chromedp.BySearch),
		//chromedp.TextContent(`.fusion-text h1`, res, chromedp.BySearch),
		chromedp.Title(res),
	}
}
