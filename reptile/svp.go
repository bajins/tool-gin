package reptile

import (
	"encoding/base64"
	"encoding/json"
	"log"
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
	urlRegex     = regexp.MustCompile(`(?:(?:https?|ftp)://)?(?:(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,63}|\[(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}|(?:[0-9a-fA-F]{1,4}:){1,7}:|::|localhost|\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(?::\d+)?(?:[/?#]\S*)?`)
	httpsRegex   = regexp.MustCompile("https?:/+(.*)")
	detailsRegex = regexp.MustCompile("(?s)<details>(.*?)</details>")

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
			var data map[string]interface{}
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
			svpMapCache.Store(typ, result)
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
	//getSvpGoroutine(&wg, 1, getSvpAbshareGit)
	//getSvpGoroutine(&wg, 2, getSvpAbshareDP)
	//getSvpGoroutine(&wg, 3, func() string {
	//	return getSvpAbshareDP1(false)
	//})
	// 密钥 (Base64)
	/*base64Key := "plr4EY25bk1HbC6a+W76TQ=="
	getSvpGoroutine(&wg, 4, func() string {
		return getSvpYse("https://api.v2rayse.com/api/live", base64Key)
	})
	getSvpGoroutine(&wg, 5, func() string {
		return getSvpYse("https://api.v2rayse.com/api/batch", base64Key)
	})*/
	getSvpGoroutine(&wg, 6, func() string {
		return getSvpYse1("https://v2rayse.com/live-node")
	})
	getSvpGoroutine(&wg, 7, func() string {
		return getSvpYse1("https://v2rayse.com/free-node")
	})
	getSvpGoroutine(&wg, 8, func() string {
		return getSvpAlvin("https://raw.githubusercontent.com/wiki/Alvin9999/new-pac/ss%E5%85%8D%E8%B4%B9%E8%B4%A6%E5%8F%B7.md")
	})
	getSvpGoroutine(&wg, 9, func() string {
		return getSvpAlvin("https://raw.githubusercontent.com/wiki/Alvin9999/new-pac/v2ray%E5%85%8D%E8%B4%B9%E8%B4%A6%E5%8F%B7.md")
	})
	getSvpGoroutine(&wg, 10, getSvpYouneed)
	// 等待所有协程完成
	wg.Wait()

	// 合并结果
	joiner := utils.NewStringJoiner("\n")
	svpMapCache.Range(func(key, value interface{}) bool {
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
