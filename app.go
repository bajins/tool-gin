package main

import (
	"embed"
	"flag"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"tool-gin/utils"
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

func init() {
	CreateTmpFiles("pyutils")
}

type embedFileSystem struct {
	http.FileSystem
}

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

	router := gin.Default()

	// 将全局中间件附加到路由器
	router.Use(FilterNoCache)
	//router.Use(Cors())
	//router.Use(Authorize())

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
	err := router.Run(Port())
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
