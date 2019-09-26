package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"tool-gin/reptile"
	"tool-gin/result"
	"tool-gin/utils"
)

// 首页
func WebRoot(c *gin.Context) {
	// 301重定向
	//c.Redirect(http.StatusMovedPermanently, "/static")
	// 返回HTML页面
	//c.HTML(http.StatusOK, "index.html", nil)
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

// 获取系统信息
func SystemInfo(c *gin.Context) {
	data := make(map[string]interface{}, 0)
	data["Version"] = utils.ToUpper(runtime.Version())
	data["cpu"] = runtime.NumCPU()
	memStatus := runtime.MemStats{}
	// 查看内存申请和分配统计信息
	runtime.ReadMemStats(&memStatus)
	// 申请的内存
	data["Mallocs"] = memStatus.Mallocs
	// 释放的内存次数
	data["Frees"] = memStatus.Frees
	// 获取当前函数或者上层函数的标识号、文件名、调用方法在当前文件中的行号
	//pc,file,line,ok := runtime.Caller(0)
	// 获取当前进程执行的cgo调用次数
	data["NumCgoCall"] = runtime.NumCgoCall()
	// 获取当前存在的go协程数
	data["NumGoroutine"] = runtime.NumGoroutine()

	c.JSON(http.StatusOK, result.Success("获取系统信息成功", data))
}

// 获取key
func GetKey(c *gin.Context) {
	// GET 获取参数内容，没有则返回空字符串
	//company := c.Query("company")
	// POST 获取的所有参数内容的类型都是 string
	company := c.PostForm("company")

	if utils.IsStringEmpty(company) {
		c.JSON(http.StatusOK, result.Error(300, "请选择公司"))
		return
	}
	app := c.PostForm("app")
	if utils.IsStringEmpty(app) {
		c.JSON(http.StatusOK, result.Error(300, "请选择产品"))
		return
	}
	version := c.PostForm("version")
	if utils.IsStringEmpty(version) {
		c.JSON(http.StatusOK, result.Error(300, "请选择版本"))
		return
	}
	// 获取当前绝对路径
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, result.SystemError())
		return
	}
	if company == "netsarang" {
		out, err := utils.ExecutePython(path.Join(dir, "pyutils", "xshell_key.py"), app, version)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusOK, result.SystemError())
			return
		}
		c.JSON(http.StatusOK, result.Success("获取key成功", map[string]string{"key": out}))

	} else if company == "mobatek" {

		_, err := utils.ExecutePython(path.Join(dir, "pyutils", "moba_xterm_Keygen.py"), utils.OsPath(), version)
		if err != nil {
			c.JSON(http.StatusOK, result.SystemError())
			return
		}
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", "Custom.mxtpro"))
		//c.Writer.Header().Set("Content-Type", "application/octet-stream")
		//c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", "Custom.mxtpro"))
		c.FileAttachment(utils.OsPath()+"/Custom.mxtpro", "Custom.mxtpro")

	} else if company == "torchsoft" {
		out, err := utils.ExecutePython(path.Join(dir, "pyutils", "reg_workshop_keygen.py"), version)
		if err != nil {
			c.JSON(http.StatusOK, result.SystemError())
			return
		}
		c.JSON(http.StatusOK, result.Success("获取key成功", map[string]string{"key": out}))
	}

}

// 文件上传请求
func Upload(c *gin.Context) {
	// 拿到上传的文件的信息
	file, header, err := c.Request.FormFile("upload")
	filename := header.Filename
	fmt.Println(header.Filename)
	out, err := os.Create("./tmp/" + filename + ".png")
	if err != nil {
		log.Println(err)
	}
	defer out.Close()
	// 拷贝上传的文件信息到新建的out文件中
	_, err = io.Copy(out, file)
	if err != nil {
		log.Println(err)
	}
}

// 文件下载请求
func Dowload(c *gin.Context) {
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

// 获取NetSarang下载url
func GetNetSarangDownloadUrl(c *gin.Context) {
	// POST 获取的所有参数内容的类型都是 string
	app := c.PostForm("app")
	if utils.IsStringEmpty(app) {
		c.JSON(http.StatusOK, result.Error(300, "请选择产品"))
		return
	}
	version := c.PostForm("version")
	if utils.IsStringEmpty(version) {
		c.JSON(http.StatusOK, result.Error(300, "请选择版本"))
		return
	}
	url, err := reptile.DownloadNetsarang(app)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, result.SystemError())
		return
	}
	c.JSON(http.StatusOK, result.Success("获取"+app+"成功", map[string]string{"url": url}))
}

// NGINX格式化代码页面
func NginxFormatIndex(c *gin.Context) {
	// 301重定向
	//c.Redirect(http.StatusMovedPermanently, "/static")
	// 返回HTML页面
	//c.HTML(http.StatusOK, "index.html", nil)
	c.HTML(http.StatusOK, "nginx-format.html", gin.H{})
}

// 格式化nginx配置代码
func NginxFormatPython(c *gin.Context) {
	// GET 获取参数内容，没有则返回空字符串
	//code := c.Query("code")
	// POST 获取的所有参数内容的类型都是 string
	code := c.PostForm("code")

	if utils.IsStringEmpty(code) {
		c.JSON(http.StatusOK, result.Error(300, "请输入配置代码"))
		return
	}
	// 获取当前绝对路径
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, result.SystemError())
		return
	}
	out, err := utils.ExecutePython(path.Join(dir, "pyutils", "nginxfmt.py"), code)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, result.SystemError())
		return
	}
	res := make(map[string]string)
	res["contents"] = out
	c.JSON(http.StatusOK, result.Success("请求成功", res))

}
