package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"tool-gin/reptile"
	"tool-gin/utils"
)

// WebRoot 首页
func WebRoot(c *gin.Context) {
	// 301重定向
	//c.Redirect(http.StatusMovedPermanently, "/static")
	// 返回HTML页面
	//c.HTML(http.StatusOK, "index.html", nil)
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

// SystemInfo 获取系统信息
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

	SuccessJSON(c, "获取系统信息成功", data)
}

// GetKey 获取key
func GetKey(c *gin.Context) {
	// GET 获取参数内容，没有则返回空字符串
	//company := c.Query("company")
	// POST 获取的所有参数内容的类型都是 string
	company := c.PostForm("company")

	if utils.IsStringEmpty(company) {
		ErrorJSON(c, 300, "请选择公司")
		return
	}
	app := c.PostForm("app")
	if utils.IsStringEmpty(app) {
		ErrorJSON(c, 300, "请选择产品")
		return
	}
	version := c.PostForm("version")
	if utils.IsStringEmpty(version) {
		ErrorJSON(c, 300, "请选择版本")
		return
	}
	dir := TempDirPath + string(filepath.Separator)
	if company == "netsarang" {
		switch strings.ToLower(app) {
		case "xshell":
			app = "Xshell"
		case "xftp":
			app = "Xftp"
		case "xlpd":
			app = "Xlpd"
		case "xmanager":
			app = "Xmanager"
		case "xshellplus":
			app = "Xshell Plus"
		case "powersuite":
			app = "Xmanager"
		}
		out, err := utils.ExecutePython(dir+"xshell_key.py", app, version)
		ExecuteScriptError(err)
		if err != nil {
			log.Println(err)
			ErrorJSON(c, http.StatusInternalServerError, "系统错误！")
			return
		}
		SuccessJSON(c, "获取key成功", map[string]string{"key": out})

	} else if company == "mobatek" {
		curr, err := utils.OsPath()
		if err != nil {
			SystemErrorJSON(c, http.StatusInternalServerError, "系统错误！")
			return
		}
		_, err = utils.ExecutePython(dir+"moba_xterm_Keygen.py", curr, version)
		ExecuteScriptError(err)
		if err != nil {
			SystemErrorJSON(c, http.StatusInternalServerError, "系统错误！")
			return
		}
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename=\"Custom.mxtpro\"")
		//c.Writer.Header().Set("Content-Type", "application/octet-stream")
		//c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", "Custom.mxtpro"))

		c.FileAttachment(filepath.Join(curr, "Custom.mxtpro"), "Custom.mxtpro")

	} else if company == "torchsoft" {
		out, err := utils.ExecutePython(dir+"reg_workshop_keygen.py", version)
		ExecuteScriptError(err)
		if err != nil {
			ErrorJSON(c, http.StatusInternalServerError, "系统错误！")
			return
		}
		SuccessJSON(c, "获取key成功", map[string]string{"key": out})
	}

}

// ExecuteScriptError 脚本执行错误处理
func ExecuteScriptError(err error) {
	// 如果命令执行错误
	if err != nil && strings.Contains(err.Error(), "exit status 1") {
		p := TempDirPath + string(filepath.Separator) + "requirements.txt"
		_, err := utils.Execute("pip", "install", "-r", p)
		if err != nil {
			return
		}
	}
}

// Upload 文件上传请求
func Upload(c *gin.Context) {
	// 拿到上传的文件的信息
	file, header, err := c.Request.FormFile("upload")
	filename := header.Filename
	log.Println(header.Filename)
	out, err := os.Create("./tmp/" + filename + ".png")
	if err != nil {
		log.Println(err)
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Println(err)
		}
	}(out)
	// 拷贝上传的文件信息到新建的out文件中
	_, err = io.Copy(out, file)
	if err != nil {
		log.Println(err)
	}
}

// Download 文件下载请求
func Download(c *gin.Context) {
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

// GetNetSarangDownloadUrl 获取NetSarang下载url
func GetNetSarangDownloadUrl(c *gin.Context) {
	// POST 获取的所有参数内容的类型都是 string
	app := c.PostForm("app")
	if utils.IsStringEmpty(app) {
		ErrorJSON(c, 300, "请选择产品")
		return
	}
	version := c.PostForm("version")
	if utils.IsStringEmpty(version) {
		ErrorJSON(c, 300, "请选择版本")
		return
	}
	url, err := reptile.NetsarangGetInfo(app)
	if err != nil {
		ErrorJSON(c, http.StatusInternalServerError, "系统错误！")
		return
	}
	SuccessJSON(c, "获取"+app+"成功", map[string]string{"url": url})
}

// NginxFormatIndex NGINX格式化代码页面
func NginxFormatIndex(c *gin.Context) {
	// 301重定向
	//c.Redirect(http.StatusMovedPermanently, "/static")
	// 返回HTML页面
	//c.HTML(http.StatusOK, "index.html", nil)
	c.HTML(http.StatusOK, "nginx-format.html", gin.H{})
}

// NginxFormatPython 格式化nginx配置代码
func NginxFormatPython(c *gin.Context) {
	// GET 获取参数内容，没有则返回空字符串
	//code := c.Query("code")
	// POST 获取的所有参数内容的类型都是 string
	code := c.PostForm("code")

	if utils.IsStringEmpty(code) {
		ErrorJSON(c, 300, "请输入配置代码")
		return
	}
	out, err := utils.ExecutePython(TempDirPath+string(filepath.Separator)+"nginxfmt.py", code)
	if err != nil {
		log.Println(err)
		ErrorJSON(c, http.StatusInternalServerError, "系统错误！")
		return
	}
	res := make(map[string]string)
	res["contents"] = out
	SuccessJSON(c, "请求成功", res)

}

// GetNavicatDownloadUrl 获取navicat下载地址
func GetNavicatDownloadUrl(c *gin.Context) {
	location, isExist := c.GetQuery("location")
	if location == "" || !isExist {
		location = c.DefaultPostForm("location", "1")
	}
	product, isExist := c.GetQuery("product")
	if product == "" || !isExist {
		product = c.DefaultPostForm("product", "navicat_premium_cs_x64.exe")
	}

	// POST 获取的所有参数内容的类型都是 string
	params := map[string]string{
		"product":    product,
		"location":   location,
		"support":    "",
		"linux_dist": "",
	}
	url := "https://www.navicat.com.cn/includes/Navicat/direct_download.php"
	result, err := utils.HttpReadBodyJsonMap(http.MethodPost, url, utils.ContentTypeAXWFU, params, nil)

	if result == nil || err != nil {
		ErrorJSON(c, http.StatusInternalServerError, "系统错误！")
		return
	}
	SuccessJSON(c, "获取下载地址成功", map[string]string{"url": result["download_link"].(string)})
}

func GetSvp(c *gin.Context) {
	defer func() { // 捕获panic
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
			c.String(http.StatusOK, r.(string))
		}
	}()
	c.String(http.StatusOK, reptile.GetSvpAll())
}
