package main

import (
	"embed"
	"flag"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"tool-gin/utils"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// 常量
const (
	// TokenSalt 可自定义盐值
	TokenSalt = "default_salt"
)

// 内嵌资源目录指令
//
//go:embed static pyutils/*[^.go]
var local embed.FS

// init函数用于初始化应用程序，将在程序启动时自动执行
func init() {
	// 这里调用了CreateTmpFiles函数，目的是在程序运行前创建必要的临时文件目录
	// 参数"pyutils"指定了创建的临时文件目录的名称
	CreateTmpFiles("pyutils")
}

type embedFileSystem struct {
	http.FileSystem
}

// Exists 检查给定路径的文件或目录是否存在。
// 该方法通过尝试打开文件来判断路径是否存在，如果能够成功打开，则认为路径存在。
// 此方法适用于嵌入式文件系统，允许程序检查资源是否可用。
// 参数:
//
//	prefix: 资源的前缀，用于区分不同的文件系统或资源集。
//	path: 要检查的文件或目录的路径。
//
// 返回值:
//
//	bool: 如果文件或目录存在，则返回true；否则返回false。
//
// 注意:
//   - 此方法依赖于e.Open方法来实际检查路径。
//   - 不存在的路径将返回false，而实际上可能是因为其他原因（如权限问题）导致的打开失败。
func (e embedFileSystem) Exists(prefix, path string) bool {
	_, err := e.Open(path)
	if err != nil {
		return false
	}
	return true
}

// EmbedFolder embed.FS转换为http.FileSystem https://github.com/gin-contrib/static/issues/19
func EmbedFolder(fsEmbed embed.FS, targetPath string) static.ServeFileSystem {
	//http.FS(os.DirFS(targetPath))
	fsys, err := fs.Sub(fsEmbed, targetPath) // 获取目录下的文件
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}

// EmbedDir 返回一个 ServeFileSystem 接口，用于提供目录服务。
// 该函数通过调用 EmbedFolder 函数，将指定路径的目录嵌入到可执行文件中。
// 参数:
//
//	targetPath - 目标目录的路径，表示要嵌入的目录。
//
// 返回值:
//
//	static.ServeFileSystem - 一个接口类型，提供了访问嵌入目录中文件和子目录的功能。
func EmbedDir(targetPath string) static.ServeFileSystem {
	return EmbedFolder(local, targetPath)
}

// Authorize 认证拦截中间件
func Authorize(c *gin.Context) {
	username := c.Query("username") // 用户名
	ts := c.Query("ts")             // 时间戳
	token := c.Query("token")       // 访问令牌

	if strings.ToLower(utils.MD5(username+ts+TokenSalt)) == strings.ToLower(token) {
		// 验证通过，会继续访问下一个中间件
		c.Next()
	} else {
		// 验证不通过，不再调用后续的函数处理
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"message": "访问未授权"})
		// return可省略, 只要前面执行Abort()就可以让后面的handler函数不再执行
		return
	}
}

// FilterNoCache 禁止浏览器页面缓存
func FilterNoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	// 继续访问下一个中间件
	c.Next()
}

// Cors 处理跨域请求,支持options访问
func Cors(c *gin.Context) {

	// 它指定允许进入来源的域名、ip+端口号 。 如果值是 ‘*’ ，表示接受任意的域名请求，这个方式不推荐，
	// 主要是因为其不安全，而且因为如果浏览器的请求携带了cookie信息，会发生错误
	c.Header("Access-Control-Allow-Origin", "*")
	// 设置服务器允许浏览器发送请求都携带cookie
	c.Header("Access-Control-Allow-Credentials", "true")
	// 允许的访问方法
	c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS, DELETE, PATCH")
	// Access-Control-Max-Age 用于 CORS 相关配置的缓存
	c.Header("Access-Control-Max-Age", "3600")
	// 设置允许的请求头信息 DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization
	c.Header("Access-Control-Allow-Headers", "Token,Origin, X-Requested-With, Content-Type, Accept,mid,X-Token,AccessToken,X-CSRF-Token, Authorization")

	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")

	method := c.Request.Method
	// 放行所有OPTIONS方法
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}
	// 继续访问下一个中间件
	c.Next()
}

// Port 获取传入参数的端口，如果没传默认值为8000
func Port() (port string) {
	flag.StringVar(&port, "p", "8000", "默认端口:8000")
	flag.Parse()
	return ":" + port

	//if len(os.Args[1:]) == 0 {
	//	return ":8000"
	//}
	//return ":" + os.Args[1]
}

// 配置gin（路由、中间件）并监听运行
func run() {

	router := NewEngine()

	// 将全局中间件附加到路由器，使用中间件处理通用、横切性的关注点
	// Gin 自带的 panic 恢复中间件
	router.Use(gin.Recovery())
	router.Use(FilterNoCache)
	//router.Use(Cors())
	//router.Use(Authorize())
	// 设置可信代理的 IP 地址或 CIDR 范围
	//router.TrustedPlatform = "CF-Connecting-IP" // 信任特定平台，来自请求头部信息
	//router.ForwardedByClientIP = true // 启用基于客户端 IP 的转发功能
	//err := router.SetTrustedProxies([]string{"127.0.0.1"})
	err := router.SetTrustedProxies(nil) // 禁用代理信任，直接使用 Request.RemoteAddr 作为客户端 IP，忽略所有代理头部
	if err != nil {
		log.Println(err)
	}

	// 在go:embed下必须指定模板路径
	t, _ := template.ParseFS(local, "static/html/*.html")
	router.SetHTMLTemplate(t)
	// 注册接口
	router.Any("/", WebRoot)
	router.POST("/getKey", GetKey)
	router.POST("/SystemInfo", SystemInfo)
	router.POST("/getXshellUrl", GetNetSarangDownloadUrl)
	router.Any("/nginx-format", NginxFormatIndex)
	router.POST("/nginx-format-py", NginxFormatPython)
	router.Any("/navicat", GetNavicatDownloadUrl)
	router.Any("/svp", GetSvp)

	/*
		//router.POST("/upload", UnifiedUpload)
		//router.POST("/download", UnifiedDownload)

		// 为 multipart forms 设置文件大小限制, 默认是32MB
		// 此处为左移位运算符, << 20 表示1MiB，8 << 20就是8MiB
		router.MaxMultipartMemory = 8 << 20
		router.POST("/upload", func(c *gin.Context) {
			// 单文件
			file, _ := c.FormFile("file")
			log.Println(file.Filename)

			// 上传文件至指定的完整文件路径
			dst := "/home/test" + file.Filename
			err := c.SaveUploadedFile(file, dst)
			if err != nil {
				log.Println(err)
			}
			c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
		})
	*/

	// 注册一个目录，gin 会把该目录当成一个静态的资源目录
	// 如 static 目录下有图片路径 index/logo.png , 你可以通过 GET /static/index/logo.png 访问到
	//router.Static("/static", "./static")

	router.StaticFS("/static", EmbedFolder(local, "static"))

	//router.Use(static.Serve("/", EmbedFolder(local, "static")))
	/*router.NoRoute(func (c *gin.Context) {
		log.Printf("%s doesn't exists, redirect on /", c.Request.URL.Path)
		c.Redirect(http.StatusMovedPermanently, "/")
	})*/

	//router.LoadHTMLFiles("./static/html/index.html")
	// 注册一个路径，gin 加载模板的时候会从该目录查找
	// 参数是一个匹配字符，如 templates/*/* 的意思是 模板目录有两层
	// gin 在启动时会自动把该目录的文件编译一次缓存，不用担心效率问题
	//router.LoadHTMLGlob("static/html/*") // 在go:embed下无效

	// listen and serve on 0.0.0.0:8080
	err = router.Run(Port())
	if err != nil {
		log.Fatal(err)
	}

	/*listener, err := net.Listen("tcp", "0.0.0.0"+Port())
	  if err != nil {
	  	panic(listener)
	  }
	  httpServer := &http.Server{
	  	Handler: router,
	  }
	  err = httpServer.Serve(listener)
	  if err != nil {
	  	panic(err)
	  }*/
}
