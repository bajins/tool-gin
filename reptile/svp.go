package reptile

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	cu "github.com/Davincible/chromedp-undetected"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/go-resty/resty/v2"
	"github.com/guonaihong/gout"
	"github.com/levigross/grequests/v2"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"tool-gin/utils"
)

type RequestCounter struct {
	count  int
	expiry time.Time
}

var (
	svpCache atomic.Value // 存储预热的缓存数据
	//ipTracker map[string]struct{} // 使用空结构体节省内存
	//ipMutex   sync.Mutex          // 用于保护ipTracker的互斥锁

	counter    = make(map[string]RequestCounter) // 创建一个计数器切片
	cacheMutex sync.RWMutex                      // 用于保护缓存的读写锁
	//cacheOnce sync.Map // map[string]*sync.Once
)

const expiryDuration = 3 * time.Minute

func init() {

	go utils.SchedulerIntervalsTimer(func() {
		defer func() { // 捕获panic
			if r := recover(); r != nil {
				log.Println("Recovered from panic:", r)
			}
		}()
		getSvpAll()
	}, time.Minute*20)

	go utils.SchedulerIntervalsTimer(func() { // 遍历并删除所有过期的条目
		// 加写锁，因为我们要删除 map 中的元素
		cacheMutex.Lock()
		defer cacheMutex.Unlock()

		now := time.Now()
		for key, rc := range counter {
			// time.Now().After(rc.expiry) 检查当前时间是否晚于记录的过期时间
			if now.After(rc.expiry) {
				delete(counter, key)
			}
		}
	}, 30*time.Second)
}

// getSvpGit 获取SVP
func getSvpGit() string {
	url := "https://raw.githubusercontent.com/abshare/abshare.github.io/main/README.md"
	result, err := utils.HttpReadBodyString(http.MethodGet, url, "", nil, nil)
	if err != nil {
		panic(err.Error())
	}
	// 匹配url
	re := regexp.MustCompile(`(?:(?:https?|ftp)://)?(?:(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,63}|\[(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|(?:[0-9a-fA-F]{1,4}:){1,7}:|::|localhost|\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(?::\d+)?(?:[/?#]\S*)?`)
	reu := regexp.MustCompile("https?:/+(.*)")
	matches := re.FindAllString(result, -1)

	// URL 预处理和去重
	seen := make(map[string]struct{})
	var urls []string
	for _, match := range matches {
		// 去除字符串中的空白字符
		str := strings.TrimSpace(match)
		if !reu.MatchString(str) {
			str = "http://" + str
		}
		// 若字符串未出现过，则添加到结果切片
		if _, exists := seen[str]; !exists {
			seen[str] = struct{}{} // 空结构体不占用额外空间
			urls = append(urls, str)
		}
	}

	header := map[string]string{
		"User-Agent": "v2rayN/7.12.5",
	}

	// 并发执行：为每个URL启动一个goroutine
	var wg sync.WaitGroup
	resultsChan := make(chan string, len(urls))

	for _, urlStr := range urls {
		wg.Add(1) // WaitGroup 计数器+1

		go func(u string) {
			defer wg.Done() // goroutine 结束时，计数器-1

			result, err := utils.HttpReadBodyString(http.MethodGet, urlStr, "", nil, header)
			if err != nil {
				return
			}
			// 去除字符串中的空白字符
			urlStr = strings.TrimSpace(result)
			// 检查字符串长度是否为 4 的倍数，验证是否为BASE64编码
			if len(urlStr)%4 != 0 {
				return
			}
			// 解码字符串，检查是否出错
			by, err := base64.StdEncoding.DecodeString(urlStr)
			if err != nil {
				return
			}
			// 将成功的结果发送到 channel
			resultsChan <- string(by)
		}(url)
	}
	// 等待所有任务完成，然后关闭channel
	wg.Wait()
	close(resultsChan)

	// 使用 strings.Builder 来高效地拼接字符串，避免性能损耗
	var builder strings.Builder
	for res := range resultsChan {
		if res != "" {
			if builder.Len() > 0 {
				builder.WriteString("\n")
			}
			builder.WriteString(res)
		}
	}
	return builder.String()
}

// getSvpDP 获取SVP
func getSvpDP() string {
	url := "https://tuijianvpn.com/1044"

	/*
		硬编码请求
	*/

	result, err := utils.HttpReadBodyString(http.MethodGet, url, "", nil, nil)
	if err != nil {
		panic(err.Error())
	}
	urls := ParseSvpHtml([]byte(result))
	if urls != "" {
		return urls
	}
	/*
		模拟浏览器
	*/

	sess := grequests.NewSession(nil)
	/*jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Failed to create cookie jar: %v", err)
	}
	sess.HTTPClient.Jar = jar*/
	resp, err := sess.Get(context.Background(), url, nil)
	if err != nil {
		panic(err.Error())
	}
	//log.Println(resp.String())

	// Step 2: 模拟表单提交评论
	resp, err = sess.Post(context.Background(),
		"https://tuijianvpn.com/wp-admin/admin-ajax.php",
		&grequests.RequestOptions{
			Data: map[string]string{
				"action":              "wpdAddComment",
				"wc_comment":          utils.RandomLower(5),
				"wc_name":             "csbxdh",
				"wc_email":            "hskdcbf@gmail.com",
				"wc_website":          "",
				"submit":              "发表评论",
				"wpdiscuz_unique_id":  "0_0",
				"comment_mail_notify": "comment_mail_notify",
				"wpd_comment_depth":   "1",
				"postId":              "1044",
			},
		})
	if err != nil {
		//fmt.Printf("表单提交失败: %v\n", err)
		panic(err.Error())
	}

	// Step 3: 再次 GET 请求获取内容
	resp, err = sess.Get(context.Background(), url, nil)
	if err != nil {
		//fmt.Printf("获取内容失败: %v\n", err)
		panic(err.Error())
	}
	urls = ParseSvpHtml(resp.Bytes())
	if urls != "" {
		return urls
	}

	/*
		调用浏览器请求
	*/

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
	return res
}

func getSvpDP1() string {
	url := "https://vpnea.com/mfjd.html"

	/*
		硬编码请求
	*/

	var cookies []*http.Cookie
	cookies = append(cookies, &http.Cookie{
		Name:  "92d9977e3c8736e482bb9e23ef9e1c3b__typecho_remember_author",
		Value: "csbxdh",
	})
	cookies = append(cookies, &http.Cookie{
		Name:  "92d9977e3c8736e482bb9e23ef9e1c3b__typecho_remember_mail",
		Value: "hskdcbf%40gmail.com",
	})
	s := ""
	err := gout.GET(url).
		SetCookies(cookies...).
		BindBody(&s). //解析响应body内容
		Do()          // 自动存储 Cookie
	if err != nil {
		panic(err.Error())
	}
	urls := ParseSvpHtml([]byte(s))
	if urls != "" {
		return urls
	}

	/*
		模拟浏览器
	*/

	// 创建 resty 客户端
	client := resty.New()
	/*jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Failed to create cookie jar: %v", err)
	}
	client.SetCookieJar(jar)*/
	client.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36")

	// 1. 首次 GET 请求（获取 Cookie）
	firstGetResp, err := client.R().Get(url) // 替换为目标 URL
	if err != nil {
		panic(fmt.Sprintf("首次 GET 请求失败: %v", err))
	}
	//fmt.Printf("首次 GET 响应状态: %d\n", firstGetResp.StatusCode())
	//fmt.Println("页面内容:", firstGetResp.String())

	// 解析HTML
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(firstGetResp.Body()))
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

	// 2. 模拟表单提交（自动携带 Cookie）
	_, err = client.R().
		//SetDebug(true).
		//SetCookies(firstGetResp.Cookies()).
		SetHeader("referer", url).
		SetFormData(map[string]string{ // 设置表单数据
			"author": "csbxdh",
			"mail":   "hskdcbf@gmail.com",
			"text":   utils.RandomLower(5),
			//"parent": cdc,
			"url": "",
			"_":   "2ec89e4c7a2e4a09e1a8cfea818253de",
		}).
		Post(caction + "?time=" + strconv.FormatInt(time.Now().UnixNano()/1e6, 10)) // 替换为评论提交 URL
	if err != nil {
		panic(fmt.Sprintf("提交评论失败: %v", err))
	}
	//fmt.Printf("评论提交状态: %d\n", postResp.StatusCode())
	//fmt.Printf("Cookies: %v\n", postResp.Header().Values("Set-Cookie"))

	/*var cookies []*http.Cookie
	for _, cookie := range postResp.Cookies() {
		cookies = append(cookies, &http.Cookie{
			Name:  cookie.Name,
			Value: cookie.Value,
		})
	}*/

	// 3. 再次 GET 请求验证评论
	secondGetResp, err := client.R().Get(url) // 同首次 URL
	if err != nil {
		panic(fmt.Sprintf("二次 GET 请求失败: %v", err))
	}
	//fmt.Printf("二次 GET 响应状态: %d\n", secondGetResp.StatusCode())
	//fmt.Println("页面内容:", secondGetResp.String())
	urls = ParseSvpHtml(secondGetResp.Body())
	if urls != "" {
		return urls
	}

	/*
		调用浏览器请求
	*/

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
	return res
}

// getSvpAll 获取SVP
func getSvpAll() string {

	// 创建 channel 用于接收结果
	/*ch1 := make(chan string)
	ch2 := make(chan string)
	ch3 := make(chan string)
	// 启动协程执行任务
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("捕获 panic:", r)
			}
		}()
		ch1 <- getSvpGit()
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("捕获 panic:", r)
			}
		}()
		ch2 <- getSvpDP()
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("捕获 panic:", r)
			}
		}()
		ch3 <- getSvpDP1()
	}()
	// 等待并收集结果
	result1 := <-ch1
	result2 := <-ch2
	result3 := <-ch3
	// 合并结果
	finalResult := result1 + "\n" + result2 + "\n" + result3*/

	var wg sync.WaitGroup
	// 创建一个带缓冲的channel来收集结果，大小和url数量一致
	//resultsChan := make(chan string, len(urls))
	results := make([]string, 3)
	// 启动协程执行任务
	wg.Add(len(results)) // WaitGroup 计数器数量
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("捕获 panic:", r)
			}
		}()
		defer wg.Done() // goroutine 结束时，计数器-1
		results[0] = getSvpGit()
		log.Println("getSvpGit() 结果：", strings.Count(results[0], "\n"))
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("捕获 panic:", r)
			}
		}()
		defer wg.Done() // goroutine 结束时，计数器-1
		results[1] = getSvpDP()
		log.Println("getSvpDP() 结果：", strings.Count(results[1], "\n"))
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("捕获 panic:", r)
			}
		}()
		defer wg.Done() // goroutine 结束时，计数器-1
		results[2] = getSvpDP1()
		log.Println("getSvpDP1() 结果：", strings.Count(results[2], "\n"))
	}()
	// 等待所有协程完成
	wg.Wait()
	// 关闭channel，表示所有结果都已经发送完毕
	//close(resultsChan)
	// 合并结果
	finalResult := results[0] + "\n" + results[1] + "\n" + results[2]

	if finalResult == "" || len(finalResult) == 0 {
		panic("没有获取到内容")
	}
	res := base64.StdEncoding.EncodeToString([]byte(finalResult))
	go func() {
		svpCache.Store(res)
	}()
	return res
}

func GetSvpAllHandler(clientIP string) string {
	now := time.Now()

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	counterEntry, exists := counter[clientIP]

	// 清理过期的计数（比如超过 3 分钟不请求就重置）
	if exists && now.After(counterEntry.expiry) {
		delete(counter, clientIP)
		exists = false
	}
	var finalResult string
	if !exists {
		// 第一次请求
		counter[clientIP] = RequestCounter{count: 1, expiry: now.Add(expiryDuration)}
		finalResult = svpCache.Load().(string)
	} else {
		// 第二次请求
		finalResult = getSvpAll()
		// 处理完后重置，下次再访问又是“第一次”
		delete(counter, clientIP)
	}

	// 第一次访问：检查缓存
	/*if val, ok := cache.Load(clientIP); ok {
		return val, nil
	}*/

	// 首次加载：确保只执行一次
	/*once, _ := cacheOnce.LoadOrStore(clientIP, &sync.Once{})
	once.(*sync.Once).Do(func() {
		cache.Store(clientIP, getSvpAll())
	})*/

	return finalResult
}
