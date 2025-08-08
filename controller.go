package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"tool-gin/reptile"
	"tool-gin/utils"

	"github.com/gin-gonic/gin"
)

// WebRoot 首页
func WebRoot(ctx *Context) {
	// 301重定向
	//ctx.C.Redirect(http.StatusMovedPermanently, "/static")
	// 返回HTML页面
	//ctx.C.HTML(http.StatusOK, "index.html", nil)
	ctx.C.HTML(http.StatusOK, "index.html", gin.H{})
}

// SystemInfo 获取系统信息
func SystemInfo(ctx *Context) {
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

	ctx.SuccessJSON("获取系统信息成功", data)
}

// GetKeyInfo 定义了GetKey的 JSON 结构体
type GetKeyInfo struct {
	Company string `json:"company" binding:"required"`
	App     string `json:"app" binding:"required"`
	Version string `json:"version" binding:"required"`
}

// GetKey 获取key
func GetKey(ctx *Context) {
	// GET 获取参数内容，没有则返回空字符串
	//company := ctx.C.Query("company")
	// POST 获取的所有参数内容的类型都是 string
	company := ctx.C.PostForm("company")

	/*var getKeyInfo GetKeyInfo
	if !ctx.C.BindAndValidate(&getKeyInfo) {
		return
	}*/
	if utils.IsStringEmpty(company) {
		ctx.ErrorJSON(300, "请选择公司")
		return
	}
	app := ctx.C.PostForm("app")
	if utils.IsStringEmpty(app) {
		ctx.ErrorJSON(300, "请选择产品")
		return
	}
	version := ctx.C.PostForm("version")
	if utils.IsStringEmpty(version) {
		ctx.ErrorJSON(300, "请选择版本")
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
			ctx.ErrorJSON(http.StatusInternalServerError, "系统错误！")
			return
		}
		ctx.SuccessJSON("获取key成功", map[string]string{"key": out})

	} else if company == "mobatek" {
		curr, err := utils.OsPath()
		if err != nil {
			ctx.SystemErrorJSON(ERROR, "系统错误！")
			return
		}
		_, err = utils.ExecutePython(dir+"moba_xterm_Keygen.py", curr, version)
		ExecuteScriptError(err)
		if err != nil {
			ctx.SystemErrorJSON(ERROR, "系统错误！")
			return
		}
		ctx.C.Header("Content-Type", "application/octet-stream")
		ctx.C.Header("Content-Disposition", "attachment; filename=\"Custom.mxtpro\"")
		//ctx.C.Writer.Header().Set("Content-Type", "application/octet-stream")
		//ctx.C.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", "Custom.mxtpro"))

		ctx.C.FileAttachment(filepath.Join(curr, "Custom.mxtpro"), "Custom.mxtpro")

	} else if company == "torchsoft" {
		out, err := utils.ExecutePython(dir+"reg_workshop_keygen.py", version)
		ExecuteScriptError(err)
		if err != nil {
			ctx.ErrorJSON(http.StatusInternalServerError, "系统错误！")
			return
		}
		ctx.SuccessJSON("获取key成功", map[string]string{"key": out})
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
func Upload(ctx *Context) {
	// 拿到上传的文件的信息
	file, header, err := ctx.C.Request.FormFile("upload")
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
func Download(ctx *Context) {
	response, err := http.Get(ctx.C.Request.Host + "/static/public/favicon.ico")
	if err != nil || response.StatusCode != http.StatusOK {
		ctx.C.Status(http.StatusServiceUnavailable)
		return
	}

	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="favicon.ico"`,
	}

	ctx.C.DataFromReader(http.StatusOK, response.ContentLength, response.Header.Get("Content-Type"), response.Body, extraHeaders)
}

// GetNetSarangDownloadUrl 获取NetSarang下载url
func GetNetSarangDownloadUrl(ctx *Context) {
	// POST 获取的所有参数内容的类型都是 string
	app := ctx.C.PostForm("app")
	if utils.IsStringEmpty(app) {
		ctx.ErrorJSON(300, "请选择产品")
		return
	}
	version := ctx.C.PostForm("version")
	if utils.IsStringEmpty(version) {
		ctx.ErrorJSON(300, "请选择版本")
		return
	}
	url, err := reptile.NetsarangGetInfo(app)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "系统错误！")
		return
	}
	ctx.SuccessJSON("获取"+app+"成功", map[string]string{"url": url})
}

// NginxFormatIndex NGINX格式化代码页面
func NginxFormatIndex(ctx *Context) {
	// 301重定向
	//ctx.C.Redirect(http.StatusMovedPermanently, "/static")
	// 返回HTML页面
	//ctx.C.HTML(http.StatusOK, "index.html", nil)
	ctx.C.HTML(http.StatusOK, "nginx-format.html", gin.H{})
}

// NginxFormatPython 格式化nginx配置代码
func NginxFormatPython(ctx *Context) {
	// GET 获取参数内容，没有则返回空字符串
	//code := ctx.C.Query("code")
	// POST 获取的所有参数内容的类型都是 string
	code := ctx.C.PostForm("code")

	if utils.IsStringEmpty(code) {
		ctx.ErrorJSON(300, "请输入配置代码")
		return
	}
	out, err := utils.ExecutePython(TempDirPath+string(filepath.Separator)+"nginxfmt.py", code)
	if err != nil {
		log.Println(err)
		ctx.ErrorJSON(http.StatusInternalServerError, "系统错误！")
		return
	}
	res := make(map[string]string)
	res["contents"] = out
	ctx.SuccessJSON("请求成功", res)
}

// GetNavicatDownloadUrl 获取navicat下载地址
func GetNavicatDownloadUrl(ctx *Context) {
	location, isExist := ctx.C.GetQuery("location")
	if location == "" || !isExist {
		location = ctx.C.DefaultPostForm("location", "1")
	}
	product, isExist := ctx.C.GetQuery("product")
	if product == "" || !isExist {
		product = ctx.C.DefaultPostForm("product", "navicat_premium_cs_x64.exe")
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
		ctx.ErrorJSON(http.StatusInternalServerError, "系统错误！")
		return
	}
	ctx.SuccessJSON("获取下载地址成功", map[string]string{"url": result["download_link"].(string)})
}

func GetSvp(ctx *Context) {
	defer func() { // 捕获panic
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
			ctx.C.String(http.StatusOK, r.(string))
		}
	}()
	//log.Println("GetSvp Header：", ctx.C.Request.Header)
	ctx.C.String(http.StatusOK, reptile.GetSvpAllHandler(getClientIP(ctx.C.Request)))
}

// getClientIP 尝试从请求中获取真实的客户端IP
func getClientIP(r *http.Request) string {
	// 检查 X-Forwarded-For 头，这是代理服务器常用的方式
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For 可能包含多个IP，取第一个
		return strings.Split(ip, ",")[0]
	}
	// 检查 X-Real-IP 头
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}
	// 如果都没有，则使用 RemoteAddr
	// RemoteAddr 格式为 "IP:port"，我们需要去掉端口
	ip, _, _ = strings.Cut(r.RemoteAddr, ":")
	return ip
}
