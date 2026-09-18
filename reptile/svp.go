package reptile

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"tool-gin/utils"

	"github.com/go-resty/resty/v2"
)

type RequestCounter struct {
	count  int
	expiry time.Time
}

var (
	urlRegex   = regexp.MustCompile(`(?:(?:https?|ftp)://)?(?:(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,63}|\[(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|(?:[0-9a-fA-F]{1,4}:){1,7}:|::|localhost|\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(?::\d+)?(?:[/?#]\S*)?`)
	httpsRegex = regexp.MustCompile("https?:/+(.*)")
	// 匹配协议开头的正则
	schemeReg = regexp.MustCompile(`(?i)(http|https|socks|socks5|v2rayn|hy2|hysteria2\+realm|hysteria2\+realm\+http|naive\+https|naive\+quic|socks4|vmess|ssr|ss|vless|trojan|hysteria2|tuic|wireguard|anytls|naive)://\S+`)

	vnVersion   atomic.Value // 存储预热的缓存数据
	smTime      sync.Map     // 存储时间
	svpMapCache sync.Map     // 存储预热的缓存数据
	svpCache    atomic.Value // 存储预热的缓存数据
	//ipTracker map[string]struct{} // 使用空结构体节省内存
	//ipMutex   sync.Mutex          // 用于保护ipTracker的互斥锁

	setSvpCacheTime atomic.Value // 创建一个计数器切片
	//ipCounter       = make(map[string]RequestCounter) // 创建一个计数器切片
	//cacheMutex      sync.RWMutex                      // 用于保护缓存的读写锁
	//cacheOnce sync.Map // map[string]*sync.Once
	reqSync atomic.Bool
)

const (
	expiryDuration = 3 * time.Minute
)

func init() {
	{
		response, err := resty.New().R().
			Get("https://api.github.com/repos/2dust/v2rayN/releases/latest")
		if err != nil {
			log.Println(err)
		}
		if response.StatusCode() == 200 {
			var data map[string]any
			err := json.Unmarshal(response.Body(), &data)
			if err != nil {
				log.Println(err)
			} else {
				vnVersion.Store("v2rayN/" + data["tag_name"].(string))
			}
		}
	}

	go utils.SchedulerIntervalsTimer(func() {
		defer func() { // 捕获panic
			if r := recover(); r != nil {
				log.Println("Recovered from panic:", r)
			}
		}()
		// 保活cookie
		getSvpAbshareDP1(true)
	}, time.Minute*20)

	/*go utils.SchedulerIntervalsTimer(func() { // 遍历并删除所有过期的条目
		// 加写锁，因为我们要删除 map 中的元素
		cacheMutex.Lock()
		defer cacheMutex.Unlock()

		now := time.Now()
		for key, rc := range ipCounter {
			// time.Now().After(rc.expiry) 检查当前时间是否晚于记录的过期时间
			if now.After(rc.expiry) {
				delete(ipCounter, key)
			}
		}
	}, 30*time.Second)*/
}

func getSvpGoroutine(wg *sync.WaitGroup, typ int, fun func() string) {
	wg.Add(1) // WaitGroup 计数器数量
	go func() {
		defer func() {
			wg.Done() // goroutine 结束时，计数器-1
			if r := recover(); r != nil {
				log.Println("捕获 panic:", r, string(debug.Stack()))
				//if err, ok := r.(error); ok && strings.Contains(err.Error(), "403") {
				t, b := smTime.Load(typ)
				if b {
					smTime.Store(typ, t.(time.Time).Add(time.Hour*3))
				}
				//}
			}
		}()
		result := fun()
		if result != "" && len(result) > 0 {
			var builder strings.Builder
			// 使用 Scanner 逐行读取
			scanner := bufio.NewScanner(strings.NewReader(result))
			isFirstLine := true
			for scanner.Scan() {
				line := scanner.Text()
				u, err := url.Parse(strings.TrimSpace(line))
				if err == nil {
					scheme := strings.ToLower(u.Scheme)
					// 判断 HTTP / HTTPS 是否为普通网页
					if scheme == "http" || scheme == "https" {
						// 如果有明确的 User 认证信息（如 http://user:pass@host:port）
						// 或者有非默认端口且没有复杂 path（如 http://1.1.1.1:8080#备注）
						isExplicitPort := (scheme == "http" && u.Port() != "" && u.Port() != "80") || (scheme == "https" && u.Port() != "" && u.Port() != "443")
						hasAuth := u.User != nil
						isFlatPath := u.Path == "" || u.Path == "/"

						if !isFlatPath || (!hasAuth && !isExplicitPort) {
							continue
						}
					}
				}
				// 判断如果不包含 127.0.0.1 则写入 Builder
				if !strings.Contains(line, "127.0.0.1") {
					if !isFirstLine {
						builder.WriteString("\n")
					}
					builder.WriteString(line)
					isFirstLine = false
				}
			}
			output := builder.String()
			svpMapCache.Store(typ, output)
		}
		log.Println("SVP ", typ, "结果：", strings.Count(result, "\n"))
	}()
}

// GetSvpAll 获取SVP
func GetSvpAll(id string) string {
	// 判断时间是否过期
	t := setSvpCacheTime.Load()
	if t != nil && !time.Now().After(t.(time.Time).Add(expiryDuration)) {
		return svpCache.Load().(string)
	}
	// 防止并发请求
	if reqSync.Load() {
		return svpCache.Load().(string)
	}
	reqSync.Store(true)
	defer reqSync.Store(false)

	// errgroup.Group 可返回 error；任意 goroutine 出错会取消其余
	var wg sync.WaitGroup
	// 启动协程执行任务
	getSvpGoroutine(&wg, 1, getSvpAbshareGit)
	getSvpGoroutine(&wg, 2, getSvpAbshareDP)
	getSvpGoroutine(&wg, 3, func() string {
		return getSvpAbshareDP1(false)
	})
	// 密钥 (Base64)
	/*base64Key := "plr4EY25bk1HbC6a+W76TQ=="
	getSvpGoroutine(&wg, 4, func() string {
		return getSvpYse("https://api.v2rayse.com/api/live", base64Key)
	})
	getSvpGoroutine(&wg, 5, func() string {
		return getSvpYse("https://api.v2rayse.com/api/batch", base64Key)
	})*/
	//getSvpGoroutine(&wg, 6, func() string {
	//	return getSvpYse1("https://v2rayse.com/live-node")
	//})
	getSvpGoroutine(&wg, 7, func() string {
		return getSvpYse1("https://nodebuf.com/free-node")
	})
	getSvpGoroutine(&wg, 8, func() string {
		// https://gitlab.com/zhifan999/fq
		return getSvpUrlReg("https://raw.githubusercontent.com/wiki/Alvin9999-newpac/fanqiang/ss免费账号.md")
	})
	getSvpGoroutine(&wg, 9, func() string {
		return getSvpUrlReg("https://raw.githubusercontent.com/wiki/Alvin9999-newpac/fanqiang/v2ray免费账号.md")
	})
	getSvpGoroutine(&wg, 9, func() string {
		// https://github.com/mahdibland/V2RayAggregator/tree/master/sub/splitted
		return getSvpUrlReg("https://raw.githubusercontent.com/mahdibland/V2RayAggregator/master/README.md")
	})
	getSvpGoroutine(&wg, 10, getSvpYouneed)
	// 等待所有协程完成
	wg.Wait()

	// 合并结果
	joiner := utils.NewStringJoiner("\n")
	svpMapCache.Range(func(key, value any) bool {
		joiner.Add(value)
		return true
	})
	// joiner.Add(getSvpGitAgg())

	if joiner.Empty() {
		panic("没有获取到内容")
	}
	finalResult := utils.RemoveDuplicateLines(joiner.String())
	res := base64.StdEncoding.EncodeToString([]byte(finalResult))
	svpCache.Store(res)
	setSvpCacheTime.Store(time.Now())
	return res
}

//func GetSvpAllHandler(clientIP string) string {
//	now := time.Now()
//
//	cacheMutex.Lock()
//	defer cacheMutex.Unlock()
//
//	counterEntry, exists := ipCounter[clientIP]
//
//	// 清理过期的计数（比如超过 3 分钟不请求就重置）
//	if exists && now.After(counterEntry.expiry) {
//		delete(ipCounter, clientIP)
//		exists = false
//	}
//	var finalResult string
//	if !exists { // IP不存在
//		// 第一次请求
//		ipCounter[clientIP] = RequestCounter{count: 1, expiry: now.Add(expiryDuration)}
//		finalResult = svpCache.Load().(string)
//	} else { // IP存在
//		// 第二次请求
//		finalResult = getSvpAll()
//		// 处理完后重置，下次再访问又是“第一次”
//		delete(ipCounter, clientIP)
//	}
//
//	// 第一次访问：检查缓存
//	/*if val, ok := cache.Load(clientIP); ok {
//		return val, nil
//	}*/
//
//	// 首次加载：确保只执行一次
//	/*once, _ := cacheOnce.LoadOrStore(clientIP, &sync.Once{})
//	once.(*sync.Once).Do(func() {
//		cache.Store(clientIP, getSvpAll())
//	})*/
//
//	return finalResult
//}

// FetchProtocolPattern 从 GitHub 获取最新的 Global.cs 并返回正则表达式字符串
// token 可选：若频繁调用可传入 GitHub Personal Access Token 避免 60次/小时 速率限制，不需要传 "" 即可
func FetchProtocolPattern(token string) (string, error) {
	apiURL := "https://api.github.com/repos/2dust/v2rayN/contents/v2rayN/ServiceLib/Global.cs?ref=master"

	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// GitHub API 必备请求头
	req.Header.Set("User-Agent", "Go-Protocol-Extractor")
	// 直接获取原始纯文本内容，避免解析 JSON 和 base64 解码
	req.Header.Set("Accept", "application/vnd.github.raw")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求 GitHub API 失败: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("GitHub API 响应异常 (HTTP %d): %s", resp.StatusCode, string(body))
	}

	contentBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应内容失败: %w", err)
	}

	// 从 Global.cs 文本中提取所有 ProtocolShare 相关的协议名（已去掉 ://）
	var protocols []string
	content := string(contentBytes)
	seen := make(map[string]bool)

	addProto := func(p string) {
		p = strings.TrimSpace(p)
		if p != "" && !seen[p] {
			seen[p] = true
			protocols = append(protocols, p)
		}
	}

	// 1. 匹配常量定义，如: public const string Hysteria2ProtocolShare = "hy2://";
	constRegex := regexp.MustCompile(`const\s+string\s+\w*Protocol\w*\s*=\s*"([^"]+)://"`)
	for _, m := range constRegex.FindAllStringSubmatch(content, -1) {
		if len(m) > 1 {
			addProto(m[1])
		}
	}

	// 2. 匹配 ProtocolShares 字典块内部的 "xxx://"
	dictBlockRegex := regexp.MustCompile(`ProtocolShares\s*=\s*new[^{]*\{([\s\S]*?)\};`)
	if dictBlock := dictBlockRegex.FindStringSubmatch(content); len(dictBlock) > 1 {
		itemRegex := regexp.MustCompile(`"([^"]+)://"`)
		for _, m := range itemRegex.FindAllStringSubmatch(dictBlock[1], -1) {
			if len(m) > 1 {
				addProto(m[1])
			}
		}
	}
	if len(protocols) == 0 {
		return "", fmt.Errorf("未在源码中解析到任何协议定义")
	}

	var escaped []string
	for _, p := range protocols {
		// 对 "+" 等正则特殊符号做转义（如 naive+https 转为 naive\+https）
		escaped = append(escaped, regexp.QuoteMeta(p))
	}
	// 将提取出的协议名列表组装为已转义的正则模式串
	return fmt.Sprintf("(?i)(%s)://\\S+", strings.Join(escaped, "|")), nil
}
