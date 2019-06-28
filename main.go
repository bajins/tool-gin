package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"key-go/common"
	"key-go/utils"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	// 可自定义盐值
	TokenSalt = "default_salt"
)

/**
 * 认证拦截中间件
 *
 * @author claer www.bajins.com
 * @date 2019/6/28 11:32
 */
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
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
}

// 禁止浏览器页面缓存
var FilterNoCache = func(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	// 继续访问下一个中间件
	c.Next()
}

// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 它指定允许进入来源的域名、ip+端口号 。 如果值是 ‘*’ ，表示接受任意的域名请求，这个方式不推荐，
		// 主要是因为其不安全，而且因为如果浏览器的请求携带了cookie信息，会发生错误
		c.Header("Access-Control-Allow-Origin", "*")
		// 设置服务器允许浏览器发送请求都携带cookie
		c.Header("Access-Control-Allow-Credentials", "true")
		// 允许的访问方法
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS, DELETE, PATCH")
		// Access-Control-Max-Age 用于 CORS 相关配置的缓存
		c.Header("Access-Control-Max-Age", "3600")
		// 设置允许的请求头信息
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
}

/**
 * 获取key
 *
 * @author claer www.bajins.com
 * @date 2019/6/28 15:04
 */
func getKey(c *gin.Context) {
	// GET 获取参数内容，没有则返回空字符串
	//company := c.Query("company")
	// POST 获取的所有参数内容的类型都是 string
	company := c.PostForm("company")

	if utils.IsStringEmpty(company) {
		c.JSON(http.StatusOK, common.Error(300, "请选择公司"))
		return
	}
	app := c.PostForm("app")
	if utils.IsStringEmpty(app) {
		c.JSON(http.StatusOK, common.Error(300, "请选择产品"))
		return
	}
	version := c.PostForm("version")
	if utils.IsStringEmpty(version) {
		c.JSON(http.StatusOK, common.Error(300, "请选择版本"))
		return
	}

	path := utils.ParentDirectory(utils.OsPath()) + "/pyutils"
	if company == "netsarang" {
		result, err := utils.ExecutePython(path+"/xshell_key.py", app, version)
		if err != nil {
			c.JSON(http.StatusOK, common.Error(500, "获取key系统错误"))
			return
		}
		res := make(map[string]string)
		res["key"] = result
		c.JSON(http.StatusOK, common.Success(200, "获取key成功", res))

	} else if company == "mobatek" {

		_, err := utils.ExecutePython(path+"/moba_xterm_Keygen.py", utils.OsPath(), version)
		if err != nil {
			c.JSON(http.StatusOK, common.Error(500, "获取key系统错误"))
			return
		}
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", "Custom.mxtpro"))
		//c.Writer.Header().Set("Content-Type", "application/octet-stream")
		//c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", "Custom.mxtpro"))
		c.FileAttachment(utils.OsPath()+"/Custom.mxtpro", "Custom.mxtpro")

	} else if company == "torchsoft" {
		result, err := utils.ExecutePython(path+"/reg_workshop_keygen.py", version)
		if err != nil {
			c.JSON(http.StatusOK, common.Error(500, "获取key系统错误"))
			return
		}
		res := make(map[string]string)
		res["key"] = result
		c.JSON(http.StatusOK, common.Success(200, "获取key成功", res))
	}

}

/**
 * 文件上传请求
 *
 * @author claer www.bajins.com
 * @date 2019/6/28 11:32
 */
func upload(c *gin.Context) {
	// 拿到上传的文件的信息
	file, header, err := c.Request.FormFile("upload")
	filename := header.Filename
	fmt.Println(header.Filename)
	out, err := os.Create("./tmp/" + filename + ".png")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	// 拷贝上传的文件信息到新建的out文件中
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
}

/**
 * 文件下载请求
 *
 * @author claer www.bajins.com
 * @date 2019/6/28 11:33
 */
func dowload(c *gin.Context) {
	response, err := http.Get(c.Request.Host + "/static/public/favicon.ico")
	if err != nil || response.StatusCode != http.StatusOK {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="favicon.ico"`,
	}

	c.DataFromReader(http.StatusOK, response.ContentLength, response.Header.Get("Content-Type"), response.Body, extraHeaders)
}

/**
 * 首页
 *
 * @Description
 * @author claer www.bajins.com
 * @date 2019/6/28 11:19
 */
func WebRoot(c *gin.Context) {
	// 301重定向
	//c.Redirect(http.StatusMovedPermanently, "/static")
	// 返回HTML页面
	//c.HTML(http.StatusOK, "index.html", nil)
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

/**
 * 获取传入参数的端口，如果没传默认值为8000
 *
 * @author claer www.bajins.com
 * @date 2019/6/28 15:31
 */
func Port() (port string) {
	flag.StringVar(&port, "p", ":8000", "默认端口:8000")
	flag.Parse()
	return port

	//if len(os.Args[1:]) == 0 {
	//	return ":8000"
	//}
	//return ":" + os.Args[1]
}

func main() {
	router := gin.Default()

	// 将全局中间件附加到路由器
	router.Use(FilterNoCache)
	//router.Use(Cors())
	//router.Use(Authorize())

	// 注册接口
	router.POST("/getKey", getKey)
	router.Any("/", WebRoot)

	// 注册一个目录，gin 会把该目录当成一个静态的资源目录
	// 该目录下的资源看可以按照路径访问
	// 如 static 目录下有图片路径 index/logo.png , 你可以通过 GET /static/index/logo.png 访问到
	router.Static("/static", "./static")
	//router.LoadHTMLFiles("./static/index.html")
	// 注册一个路径，gin 加载模板的时候会从该目录查找
	// 参数是一个匹配字符，如 templates/*/* 的意思是 模板目录有两层
	// gin 在启动时会自动把该目录的文件编译一次缓存，不用担心效率问题
	router.LoadHTMLGlob("templates/*")

	// listen and serve on 0.0.0.0:8080
	router.Run(Port())
}
