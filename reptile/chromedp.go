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
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"strings"
	"time"
	"tool-gin/utils"
)

// 启动，建议在主入口处调用一次即可
//
// context.Context部分不能抽离，否则会报 context canceled
func Apply(debug bool) (context.Context, context.CancelFunc) {
	// 创建缓存目录
	//dir, err := os.MkdirTemp("", "chromedp-example")
	//if err != nil {
	//	panic(err)
	//}
	//defer os.RemoveAll(dir)

	//dir, err := os.MkdirTemp("", "chromedp-example")
	//if err != nil {
	//	panic(err)
	//}
	//defer os.RemoveAll(dir)
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// 禁用GPU，不显示GUI
		chromedp.DisableGPU,
		// 取消沙盒模式
		chromedp.NoSandbox,
		// 指定浏览器分辨率
		//chromedp.WindowSize(1600, 900),
		// 设置UA，防止有些页面识别headless模式
		chromedp.UserAgent(utils.UserAgent),
		// 隐身模式启动
		chromedp.Flag("incognito", true),
		// 忽略证书错误
		chromedp.Flag("ignore-certificate-errors", true),
		// 窗口最大化
		chromedp.Flag("start-maximized", true),
		// 不加载图片, 提升速度
		chromedp.Flag("disable-images", true),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		// 禁用扩展
		chromedp.Flag("disable-extensions", true),
		// 禁止加载所有插件
		chromedp.Flag("disable-plugins", true),
		// 禁用浏览器应用
		chromedp.Flag("disable-software-rasterizer", true),
		//chromedp.Flag("remote-debugging-port","9222"),
		//chromedp.Flag("debuggerAddress","127.0.0.1:9222"),
		chromedp.Flag("user-data-dir", "./.cache"),
		//chromedp.Flag("excludeSwitches", "enable-automation"),
		// 设置用户数据目录
		//chromedp.UserDataDir(dir),
		//chromedp.ExecPath("C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe"),
	)
	if debug {
		opts = append(opts, chromedp.Flag("headless", false), chromedp.Flag("hide-scrollbars", false))
	}

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	// 自定义记录器
	ctx, cancel = chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
	// 设置超时时间
	ctx, cancel = context.WithTimeout(ctx, 3*time.Minute)
	//ctx, cancel = context.WithCancel(ctx)
	//if close {
	//	defer cancel()
	//}
	return ctx, cancel
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
			return os.WriteFile("testimonials.png", buf, 0644)
		}),
	}
	/*if err := os.WriteFile("fullScreenshot.png", buf, 0644); err != nil {
		log.Fatal("生成图片错误：", err)
	}*/
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

// 执行js
// https://github.com/chromedp/chromedp/issues/256
func EvalJS(js string) chromedp.Tasks {
	var res *runtime.RemoteObject
	return chromedp.Tasks{
		chromedp.EvaluateAsDevTools(js, &res),
		//chromedp.Evaluate(js, &res),
		chromedp.ActionFunc(func(ctx context.Context) error {
			b, err := res.MarshalJSON()
			if err != nil {
				return err
			}
			fmt.Println("result: ", string(b))
			return nil
		}),
	}
}

// see: https://intoli.com/blog/not-possible-to-block-chrome-headless/
const script = `(function(w, n, wn) {
	console.log(navigator.webdriver);

	// Pass the Webdriver Test.
	// chrome 为undefined，Firefox 为false
	//Object.defineProperty(n, 'webdriver', {
	//	get: () => undefined,
	//});
	// 通过原型删除该属性
	delete navigator.__proto__.webdriver;
	console.log(navigator.webdriver);
	
	// Pass the Plugins Length Test.
	// Overwrite the plugins property to use a custom getter.
	Object.defineProperty(n, 'plugins', {
	// This just needs to have length > 0 for the current test,
	// but we could mock the plugins too if necessary.
	get: () =>[
			{filename:'internal-pdf-viewer'},
			{filename:'adsfkjlkjhalkh'},
			{filename:'internal-nacl-plugin'}
		],
	});
	
	// Pass the Languages Test.
	// Overwrite the plugins property to use a custom getter.
	Object.defineProperty(n, 'languages', {
	get: () => ['zh-CN', 'en'],
	});

	// store the existing descriptor
	const elementDescriptor = Object.getOwnPropertyDescriptor(HTMLElement.prototype, 'offsetHeight');
	
	// redefine the property with a patched descriptor
	Object.defineProperty(HTMLDivElement.prototype, 'offsetHeight', {
	  ...elementDescriptor,
	  get: function() {
		if (this.id === 'modernizr') {
			return 1;
		}
		return elementDescriptor.get.apply(this);
	  },
	});

	['height', 'width'].forEach(property => {
	  // store the existing descriptor
	  const imageDescriptor = Object.getOwnPropertyDescriptor(HTMLImageElement.prototype, property);
	
	  // redefine the property with a patched descriptor
	  Object.defineProperty(HTMLImageElement.prototype, property, {
		...imageDescriptor,
		get: function() {
		  // return an arbitrary non-zero dimension if the image failed to load
		  if (this.complete && this.naturalHeight == 0) {
			return 20;
		  }
		  // otherwise, return the actual dimension
		  return imageDescriptor.get.apply(this);
		},
	  });
	});
	
	// Pass the Chrome Test.
	// We can mock this in as much depth as we need for the test.
	w.chrome = {
		runtime: {},
	};
	window.navigator.chrome = {
	  	runtime: {},
	};
	
	// Pass the Permissions Test.
	const originalQuery = wn.permissions.query;
	return wn.permissions.query = (parameters) => (
	parameters.name === 'notifications' ?
	  Promise.resolve({ state: Notification.permission }) :
	  originalQuery(parameters)
	);

})(window, navigator, window.navigator);`

// 反检测Headless
// https://github.com/chromedp/chromedp/issues/396
func AntiDetectionHeadless() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			identifier, err := page.AddScriptToEvaluateOnNewDocument(script).Do(ctx)
			if err != nil {
				return err
			}
			fmt.Println("identifier: ", identifier.String())
			return nil
		}),
	}
}
