package main

// 导包
import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tool-gin/reptile"
	"tool-gin/utils"
)

// 初始化函数
func init() {
	// 设置日志初始化参数
	// log.Lshortfile 简要文件路径，log.Llongfile 完整文件路径
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	// 设置项目为发布环境
	//gin.SetMode(gin.ReleaseMode)

	go utils.SchedulerFixedTicker(reptile.NetsarangDownloadAll, time.Hour*24)

	go utils.SchedulerFixedTicker(reptile.GetSvpDP, time.Minute*30)
}

// 运行主体函数
func main() {

	go run()

	// 通过WaitGroup管理两个协程，主协程等待两个子协程退出
	/*noChan := make(chan int)
	waiter := &sync.WaitGroup{}
	waiter.Add(2)
	go func(ch chan int, wt *sync.WaitGroup) {
		data := <-ch
		log.Println("receive data ", data)
		wt.Done()
	}(noChan, waiter)

	go func(ch chan int, wt *sync.WaitGroup) {
		ch <- 5
		log.Println("send data ", 5)
		wt.Done()
	}(noChan, waiter)
	waiter.Wait()*/

	// Go 通过向一个通道发送 `os.Signal` 值来进行信号通知。
	// 创建一个通道来接收这些通知（同时还创建一个用于在程序可以结束时进行通知的通道）。
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// `signal.Notify` 注册这个给定的通道用于接收特定信号。
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//signal.Notify(sigs, os.Interrupt)

	// 启用Go协程执行一个阻塞的信号接收操作。
	go func() {
		/*for s := range sigs {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:

				//os.Exit(1)
			default:
				log.Println("other", s)
			}
		}*/
		/*select {
		case sig := <-sigs:
			log.Printf("Got %s signal. Aborting...\n", sig)
			//os.Exit(1)
		}*/

		// 得到一个信号值
		sig := <-sigs
		log.Println("得到一个信号值：", sig)

		DestroyTempDir()

		// 通知程序可以退出
		done <- true
	}()

	// 程序将在这里进行等待，直到它得到了期望的信号
	// （也就是上面的 Go 协程发送的 `done` 值）然后退出。
	<-done
}
