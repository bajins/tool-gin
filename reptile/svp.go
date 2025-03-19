package reptile

import (
	"context"
	"encoding/base64"
	"fmt"
	cu "github.com/Davincible/chromedp-undetected"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
	"tool-gin/utils"
)

var (
	SvpMap = make([]string, 2)
)

// getSvp 获取SVP
func getSvp() string {
	result, err := utils.HttpReadBodyString(http.MethodGet, "https://raw.githubusercontent.com/abshare/abshare.github.io/main/README.md", "", nil, nil)
	if err != nil {
		panic(err.Error())
	}
	// 匹配url
	re := regexp.MustCompile(`(?:(?:https?|ftp)://)?(?:(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,63}|\[(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|(?:[0-9a-fA-F]{1,4}:){1,7}:|::|localhost|\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(?::\d+)?(?:[/?#]\S*)?`)
	reu := regexp.MustCompile("https?:/+(.*)")
	matches := re.FindAllString(result, -1)
	uniqueURLs := make(map[string]bool)

	var content string
	for _, match := range matches {
		// 去除字符串中的空白字符
		str := strings.TrimSpace(match)
		if !reu.MatchString(str) {
			str = "http://" + str
		}
		// 重复的不执行
		if uniqueURLs[str] {
			continue
		}
		uniqueURLs[str] = true
		result, err := utils.HttpReadBodyString(http.MethodGet, str, "", nil, nil)
		if err != nil {
			continue
		}
		// 去除字符串中的空白字符
		str = strings.TrimSpace(result)
		// 检查字符串长度是否为 4 的倍数，验证是否为BASE64编码
		if len(str)%4 != 0 {
			continue
		}
		// 解码字符串，检查是否出错
		by, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			continue
		}
		if content != "" {
			content += "\n"
		}
		content += string(by)
		log.Println("content:", content)
	}
	return content
}

// getSvpDP 获取SVP
func getSvpDP() string {
	ctx, cancel, err := cu.New(cu.NewConfig(
		// Remove this if you want to see a browser window.
		cu.WithHeadless(),
		// If the webelement is not found within 10 seconds, timeout.
		cu.WithTimeout(30*time.Second),
		cu.WithChromeFlags(
			chromedp.Flag("disable-application-cache", true), // 禁用应用缓存
			//chromedp.Flag("disk-cache-dir", ""),              // 禁用磁盘缓存，可能会导致加载缓慢
			//chromedp.Flag("no-cache", true),                  // 禁用内存缓存，可能会导致加载缓慢
			// 不加载图片, 提升速度
			chromedp.Flag("disable-images", true),
			chromedp.Flag("blink-settings", "imagesEnabled=false"),
		),
	))
	if err != nil {
		panic(err)
	}
	//ctx, cancel := reptile.Apply(false)
	defer cancel()

	url := "https://tuijianvpn.com/1044"
	// 定义变量，用来保存爬虫的数据
	var res string
	/*tags, _ := chromedp.Targets(ctx)
	log.Println("当前浏览器实例标签页数：", len(tags))
	if tags == nil {
		// 新建浏览器标签页及上下文
		//ctx, cancel := chromedp.NewContext(ctx, chromedp.WithTargetID(target.ID(target.CreateTarget(url).BrowserContextID)))
		//defer cancel()
	}*/
	// 随机字符串
	randStr := utils.RandomLower(5)
	name := utils.RandomMixed(5)

	// 监听 JavaScript 对话框事件
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch e := ev.(type) {
		//case *runtime.EventInspectRequested:
		//case *runtime.EventConsoleAPICalled:
		case *page.EventJavascriptDialogOpening:
			log.Printf("Alert 检测到: %s", e.Message)
			// 自动点击“确定”关闭弹窗
			if err := page.HandleJavaScriptDialog(true).Do(ctx); err != nil {
				log.Fatal(err)
			}
		}
	})

	err = chromedp.Run(ctx, chromedp.Tasks{
		//network.ClearBrowserCache(),    // 清除所有缓存，可能会导致加载缓慢
		network.SetCacheDisabled(true), // 禁用缓存
		// 浏览器下载行为，注意设置顺序，如果不是第一个会失败
		//page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorDeny),
		// 跳转页面
		chromedp.Navigate(url),
		// 等待 body 元素准备好 (DOMContentLoaded 事件后 body 通常就存在)
		/*chromedp.WaitReady("#wpd-editor-0_0 > div.ql-editor > p", chromedp.BySearch),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("DOMContentLoaded event triggered, DOM is ready.")
			return nil
		}),*/
		// 查找并等待可见
		//chromedp.WaitVisible("#wpd-editor-0_0 > div.ql-editor > p", chromedp.BySearch),
		// 点击元素
		chromedp.Click(`#wpd-editor-0_0 > div.ql-editor > p`, chromedp.BySearch),
		chromedp.ActionFunc(func(ctx context.Context) error {
			// 通过 JS 设置文本
			return chromedp.Evaluate(fmt.Sprintf(`
		                document.querySelector('#wpd-editor-0_0 > div.ql-editor > p').textContent = '%s';
		            `, randStr), nil).Do(ctx)
		}),
		chromedp.SendKeys(`#wc_name-0_0`, name, chromedp.BySearch),
		chromedp.SendKeys(`#wc_email-0_0`, utils.RandomMixed(5)+"@gmail.com", chromedp.BySearch),
		// 点击元素
		chromedp.Click(`#wpd-field-submit-0_0`, chromedp.BySearch),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("发表评论")
			return nil
		}),
		//chromedp.Sleep(3 * time.Second),
		// 查找并等待可见
		chromedp.WaitVisible(fmt.Sprintf(`img[alt="%s"]`, name), chromedp.BySearch),
		// 覆盖 beforeunload 事件，阻止弹窗
		/*chromedp.Evaluate(`
			window.onbeforeunload = null;
		    window.addEventListener('beforeunload', function(e) {
		        e.preventDefault();
		        e.returnValue = ''; // 兼容某些浏览器
		    });
		`, nil),*/
		//chromedp.Evaluate("location.reload(true);", nil), // true 强制刷新
		chromedp.Reload(), // 执行页面刷新
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("Refreshing page...")
			return nil
		}),
		// 查找并等待可见
		chromedp.WaitVisible(`.su-box-content.su-u-clearfix.su-u-trim pre`, chromedp.BySearch),
		//chromedp.Text(`.su-box-content.su-u-clearfix.su-u-trim pre`, &res, chromedp.BySearch),
		// JS更好的获取值，原生CSS selector和XPath不支持匹配到相同标签元素时获取第几个
		chromedp.Evaluate(`document.querySelectorAll(".su-box-content.su-u-clearfix.su-u-trim pre")[1].innerText`, &res),
	})
	log.Println("res:", res)
	return res
}

func getSvpDP1() string {
	ctx, cancel, err := cu.New(cu.NewConfig(
		// Remove this if you want to see a browser window.
		cu.WithHeadless(),
		// If the webelement is not found within 10 seconds, timeout.
		cu.WithTimeout(30*time.Second),
		cu.WithChromeFlags(
			chromedp.Flag("disable-application-cache", true), // 禁用应用缓存
			//chromedp.Flag("disk-cache-dir", ""),              // 禁用磁盘缓存，可能会导致加载缓慢
			//chromedp.Flag("no-cache", true),                  // 禁用内存缓存，可能会导致加载缓慢
			// 不加载图片, 提升速度
			chromedp.Flag("disable-images", true),
			chromedp.Flag("blink-settings", "imagesEnabled=false"),
		),
	))
	if err != nil {
		panic(err)
	}
	//ctx, cancel := reptile.Apply(false)
	defer cancel()

	url := "https://vpnea.com/mfjd.html"
	// 定义变量，用来保存爬虫的数据
	var res string
	/*tags, _ := chromedp.Targets(ctx)
	log.Println("当前浏览器实例标签页数：", len(tags))
	if tags == nil {
		// 新建浏览器标签页及上下文
		//ctx, cancel := chromedp.NewContext(ctx, chromedp.WithTargetID(target.ID(target.CreateTarget(url).BrowserContextID)))
		//defer cancel()
	}*/
	// 随机字符串
	randStr := utils.RandomLower(5)
	name := utils.RandomMixed(5)

	// 监听 JavaScript 对话框事件
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch e := ev.(type) {
		//case *runtime.EventInspectRequested:
		//case *runtime.EventConsoleAPICalled:
		case *page.EventJavascriptDialogOpening:
			log.Printf("Alert 检测到: %s", e.Message)
			// 自动点击“确定”关闭弹窗
			if err := page.HandleJavaScriptDialog(true).Do(ctx); err != nil {
				log.Fatal(err)
			}
		}
	})

	err = chromedp.Run(ctx, chromedp.Tasks{
		//network.ClearBrowserCache(),    // 清除所有缓存，可能会导致加载缓慢
		network.SetCacheDisabled(true), // 禁用缓存
		// 浏览器下载行为，注意设置顺序，如果不是第一个会失败
		//page.SetDownloadBehavior(page.SetDownloadBehaviorBehaviorDeny),
		// 跳转页面
		chromedp.Navigate(url),
		// 点击元素
		chromedp.Click(`textarea.text.joe_owo__target`, chromedp.BySearch),
		chromedp.SendKeys(`textarea.text.joe_owo__target`, randStr, chromedp.BySearch),
		chromedp.SendKeys(`input[name="author"]`, name, chromedp.BySearch),
		chromedp.SendKeys(`input[name="mail"]`, utils.RandomMixed(5)+"@gmail.com", chromedp.BySearch),
		// 点击元素
		chromedp.Click(`#respond-page-2 > form > div.foot > div.submit > button`, chromedp.BySearch),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Println("发表评论")
			return nil
		}),
		// 查找并等待可见
		chromedp.WaitVisible(`#Joe > div.joe_container > div > div.joe_detail > article code`, chromedp.BySearch),
		// JS更好的获取值，原生CSS selector和XPath不支持匹配到相同标签元素时获取第几个
		chromedp.Evaluate(`document.querySelectorAll(".joe_container code")[1].innerText`, &res),
	})
	log.Println("res:", res)
	return res
}

func GetSvpDP() {
	SvpMap[0] = getSvp()
	SvpMap[1] = getSvpDP()
}

// GetSvpAll 获取SVP
func GetSvpAll() string {
	/*
		// 创建 channel 用于接收结果
		ch1 := make(chan string)
		ch2 := make(chan string)
		// 启动协程执行任务
		go func() {
			ch1 <- getSvp()
		}()
		go func() {
			ch2 <- getSvpDP()
		}()
		// 等待并收集结果
		result1 := <-ch1
		result2 := <-ch2
		// 合并结果
		finalResult := result1 + "\n" + result2
	*/
	/*var wg sync.WaitGroup
	results := make([]string, 3)
	// 启动协程执行任务
	wg.Add(len(results))
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("捕获 panic:", r)
			}
		}()
		defer wg.Done()
		results[0] = getSvp()
		log.Println("getSvp() 结果：", len(results[0]))
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("捕获 panic:", r)
			}
		}()
		defer wg.Done()
		results[1] = getSvpDP()
		log.Println("getSvpDP() 结果：", len(results[1]))
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("捕获 panic:", r)
			}
		}()
		defer wg.Done()
		results[2] = getSvpDP1()
		log.Println("getSvpDP1() 结果：", len(results[2]))
	}()
	// 等待所有协程完成
	wg.Wait()
	// 合并结果
	finalResult := results[0] + "\n" + results[1] + "\n" + results[2]*/
	res := base64.StdEncoding.EncodeToString([]byte(getSvp() + "\n" + SvpMap[0] + "\n" + SvpMap[1]))
	if res == "" || len(res) == 0 {
		panic("没有获取到内容")
	}
	return res
}
