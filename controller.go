package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"tool-gin/reptile"
	"tool-gin/utils"

	"github.com/gin-gonic/gin"
)

const (
	// 上传文件最终保存的目录
	uploadDir = "./uploads"
	// 临时文件/分片存放的目录
	tempDir = "./tmp"
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

// UnifiedUpload 是一个统一的文件上传处理器
// 它可以自动处理普通上传和分片上传
func UnifiedUpload(ctx *Context) {
	// 确保上传和临时目录存在
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Println("创建上传目录失败:", err)
		ctx.ErrorJSON(http.StatusInternalServerError, "服务器内部错误")
		return
	}
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		log.Println("创建临时目录失败:", err)
		ctx.ErrorJSON(http.StatusInternalServerError, "服务器内部错误")
		return
	}

	// 检查是否为分片上传
	chunkNumber := ctx.C.PostForm("chunkNumber")
	totalChunks := ctx.C.PostForm("totalChunks")

	if chunkNumber != "" && totalChunks != "" {
		// --- 分片上传逻辑 ---
		handleChunkedUpload(ctx, chunkNumber, totalChunks)
	} else {
		// --- 普通单文件上传逻辑 ---
		handleSimpleUpload(ctx)
	}
}

// handleSimpleUpload 处理标准的、非分片的单文件上传，适用于小文件
func handleSimpleUpload(ctx *Context) {
	// 1. 获取上传文件
	// 使用 "file" 作为统一的表单字段名
	file, header, err := ctx.C.Request.FormFile("file")
	if err != nil {
		log.Println("获取文件失败:", err)
		ctx.ErrorJSON(http.StatusBadRequest, "获取上传文件失败")
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Println("关闭文件失败:", err)
		}
	}(file)

	filename := header.Filename
	log.Println("接收到普通上传文件:", filename)

	// 2. 创建临时文件，防止上传中断导致文件不完整
	tmpFile, err := os.CreateTemp(tempDir, "upload-*.tmp")
	if err != nil {
		log.Println("创建临时文件失败:", err)
		ctx.ErrorJSON(http.StatusInternalServerError, "服务器内部错误")
		return
	}
	// 使用 defer 确保临时文件最终会被关闭和清理
	defer func() {
		err := tmpFile.Close()
		if err != nil {
			return
		}
		// 如果请求处理不成功，则删除临时文件
		if ctx.C.Writer.Header().Get("Content-Type") != "application/json" { // 这是一个简化的判断方式
			log.Println("上传未成功，删除临时文件:", tmpFile.Name())
			err := os.Remove(tmpFile.Name())
			if err != nil {
				return
			}
		}
	}()

	// 3. 将上传内容拷贝到临时文件
	// 使用 io.Copy 替代手动循环读写，更简洁高效
	if _, err := io.Copy(tmpFile, file); err != nil {
		log.Println("写入临时文件失败:", err)
		ctx.ErrorJSON(http.StatusInternalServerError, "文件保存失败")
		return
	}

	// 4. 将临时文件移动到最终位置
	// 使用 Close() 确保所有内容都已刷到磁盘
	if err := tmpFile.Close(); err != nil {
		log.Println("关闭临时文件失败:", err)
		ctx.ErrorJSON(http.StatusInternalServerError, "文件保存失败")
		return
	}

	finalPath := filepath.Join(uploadDir, filename)
	if err := os.Rename(tmpFile.Name(), finalPath); err != nil {
		log.Println("重命名文件失败:", err)
		ctx.ErrorJSON(http.StatusInternalServerError, "文件保存失败")
		return
	}

	log.Println("文件已成功保存到:", finalPath)
	ctx.SuccessJSON("文件上传成功", nil)
}

// handleChunkedUpload 处理分片上传的逻辑，适用于大文件、断点续传和多线程上传
// 它会自动处理分片上传的元数据，并确保所有分片都保存到临时目录中
// chunkNumber: 当前分片的序号（从1开始或从0开始，只要保持一致即可）。
// totalChunks: 文件被分成的总片数。
// filename: 原始文件的完整名称（推荐，用于标识分片属于哪个文件）。
func handleChunkedUpload(ctx *Context, chunkNumber, totalChunks string) {
	// 解析 multipart form
	// 设置一个较大的内存限制（例如32MB），超出部分将写入临时文件
	if err := ctx.C.Request.ParseMultipartForm(32 << 20); err != nil {
		ctx.ErrorJSON(http.StatusBadRequest, "解析表单失败: "+err.Error())
		return
	}
	// 1. 获取文件分片
	file, header, err := ctx.C.Request.FormFile("file")
	if err != nil {
		ctx.ErrorJSON(http.StatusBadRequest, "获取文件分片失败: "+err.Error())
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Println("关闭文件失败:", err)
		}
	}(file)

	// 2. 获取其他元数据
	// 前端通常会发送一个唯一标识符来区分不同文件的上传
	// 这里我们简化处理，直接使用原始文件名作为标识
	originalFilename := ctx.C.PostForm("filename")
	if originalFilename == "" {
		originalFilename = header.Filename // 作为备用
	}

	log.Printf("接收到分片: %s, 文件名: %s, 分片号: %s/%s", header.Filename, originalFilename, chunkNumber, totalChunks)

	// 3. 创建用于存放该文件所有分片的临时目录
	chunksDir := filepath.Join(tempDir, originalFilename+"_chunks")
	if err := os.MkdirAll(chunksDir, os.ModePerm); err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "创建分片目录失败")
		return
	}

	// 4. 保存分片文件
	chunkFilePath := filepath.Join(chunksDir, chunkNumber)
	chunkFile, err := os.Create(chunkFilePath)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "创建分片文件失败")
		return
	}
	defer func(chunkFile *os.File) {
		err := chunkFile.Close()
		if err != nil {
			log.Println("关闭分片文件失败:", err)
		}
	}(chunkFile)

	if _, err = io.Copy(chunkFile, file); err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "写入分片数据失败")
		return
	}

	// 5. 检查所有分片是否已上传完毕
	if isUploadComplete(chunksDir, totalChunks) {
		log.Println("所有分片已上传，开始合并:", originalFilename)
		// 合并分片
		if err := mergeChunks(chunksDir, originalFilename); err != nil {
			ctx.ErrorJSON(http.StatusInternalServerError, "合并分片失败: "+err.Error())
			return
		}
		// 清理临时分片目录
		log.Println("合并成功，清理分片目录:", chunksDir)
		_ = os.RemoveAll(chunksDir)
		ctx.SuccessJSON("文件上传成功", map[string]string{"status": "merged"})
	} else {
		// 如果未完成，则只返回分片上传成功的消息
		ctx.SuccessJSON("分片上传成功", map[string]string{"status": "chunk_uploaded"})
	}
}

// isUploadComplete 检查所有分片是否已上传
func isUploadComplete(chunksDir, totalChunksStr string) bool {
	expectedChunks, err := strconv.Atoi(totalChunksStr)
	if err != nil {
		log.Println("无效的总分片数:", totalChunksStr)
		return false
	}

	files, err := os.ReadDir(chunksDir)
	if err != nil {
		log.Println("读取分片目录失败:", err)
		return false
	}

	return len(files) == expectedChunks
}

// mergeChunks 合并所有分片到一个最终文件中
func mergeChunks(chunksDir, finalFilename string) error {
	// 1. 创建最终文件
	finalPath := filepath.Join(uploadDir, finalFilename)
	finalFile, err := os.Create(finalPath)
	if err != nil {
		return fmt.Errorf("创建最终文件失败: %w", err)
	}
	defer func(finalFile *os.File) {
		err := finalFile.Close()
		if err != nil {
			log.Println("关闭最终文件失败:", err)
		}
	}(finalFile)

	// 2. 获取分片列表并按数字顺序排序
	chunkEntries, err := os.ReadDir(chunksDir)
	if err != nil {
		return fmt.Errorf("读取分片目录失败: %w", err)
	}

	sort.Slice(chunkEntries, func(i, j int) bool {
		// 文件名就是分片号，直接转为整数比较
		numI, _ := strconv.Atoi(chunkEntries[i].Name())
		numJ, _ := strconv.Atoi(chunkEntries[j].Name())
		return numI < numJ
	})

	// 3. 逐个追加分片内容到最终文件
	for _, entry := range chunkEntries {
		chunkPath := filepath.Join(chunksDir, entry.Name())
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			return fmt.Errorf("打开分片 %s 失败: %w", entry.Name(), err)
		}

		_, err = io.Copy(finalFile, chunkFile)
		_ = chunkFile.Close()
		if err != nil {
			return fmt.Errorf("合并分片 %s 失败: %w", entry.Name(), err)
		}
	}

	log.Println("所有分片成功合并到:", finalPath)
	return nil
}

// UnifiedDownload 是一个统一的文件下载处理器。
// 它能自动处理普通下载和基于 Range 的分片下载。
func UnifiedDownload(ctx *Context) {
	// 1. 从 URL 参数中获取文件名
	// 例如，对于 /download/my-video.mp4，filename 会是 "my-video.mp4"
	filename := ctx.C.Param("filename")
	// 1. 从 URL 中解析出文件名
	// 例如，从 "/download/my-video.mp4" 中提取 "my-video.mp4"
	// 我们假设文件名是 URL 路径的最后一部分
	//filename := filepath.Base(ctx.C.Request.URL.Path)
	// 基本的安全检查，防止空的或非法的路径段
	if filename == "" || filename == "." || filename == ".." {
		ctx.ErrorJSON(http.StatusBadRequest, "无效的文件名")
		return
	}

	// 2. 构建文件的完整、安全路径
	// 使用 filepath.Join 可以防止路径遍历攻击 (e.g., /download/../../etc/passwd)
	// filepath.Base 已经移除了目录部分，这里是双重保险
	filePath := filepath.Join(uploadDir, filename)

	// 3. 检查文件是否存在且不是一个目录
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		log.Printf("文件未找到: %s", filePath)
		ctx.ErrorJSON(http.StatusNotFound, "文件未找到 (404 Not Found)")
		return
	}
	if err != nil {
		log.Printf("检查文件时出错: %v", err)
		ctx.ErrorJSON(http.StatusInternalServerError, "无效的文件名")
		return
	}
	if fileInfo.IsDir() {
		log.Printf("请求的是一个目录，非法访问: %s", filePath)
		ctx.ErrorJSON(http.StatusBadRequest, "无效的请求，不能下载目录")
		return
	}

	// 5. 使用 http.ServeFile 提供文件服务
	// 这是最关键的一步。ServeFile 会自动处理：
	// - 设置正确的 Content-Type (基于文件扩展名)
	// - 设置 Content-Length
	// - 处理 If-Modified-Since 请求 (实现浏览器缓存)
	// - 设置 Content-Disposition: attachment; filename="...". "attachment" 表示附件形式，浏览器应提示用户保存
	// - 自动检测并响应 Range 请求，实现分片下载！
	log.Printf("开始提供文件下载: %s (Range: %s)", filePath, ctx.C.Request.Header.Get("Range"))
	ctx.C.FileAttachment(filePath, filename)

	/*
		// 使用 bytes.Buffer 在内存中构建文件，它实现了 io.Reader
		buffer := new(bytes.Buffer)
		file, err := os.ReadFile(filePath)
		if err != nil {
			ctx.ErrorJSON(http.StatusInternalServerError, "读取文件失败")
			return
		}
		buffer.Write(file)
		ctype := mime.TypeByExtension(filepath.Ext(filename))
		if ctype == "" {
			// read a chunk to decide between utf-8 text and binary
			var buf [512]byte
			n, _ := io.ReadFull(buffer, buf[:])
			ctype = http.DetectContentType(buf[:n])
		}
		extraHeaders := map[string]string{
			// 关键：设置这个头来触发浏览器下载
			"Content-Disposition": `attachment; filename="` + filename + `"`,
		}
		// 流式处理方法，从 Reader 中读取数据块，然后直接写入到 HTTP 响应中，而不会将所有数据一次性加载到内存里
		ctx.C.DataFromReader(http.StatusOK, fileInfo.Size(), ctype, buffer, extraHeaders)
	*/
}

// MultiThreadDownloadByUrl 多线程分片下载远程文件
func MultiThreadDownloadByUrl(ctx *Context) {
	// 获取要下载的文件URL
	fileURL := ctx.C.Query("url")
	if utils.IsStringEmpty(fileURL) {
		ctx.ErrorJSON(http.StatusBadRequest, "请提供文件URL")
		return
	}

	// 获取文件信息
	resp, err := http.Head(fileURL)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "无法获取文件信息")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("关闭响应体失败:", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		ctx.ErrorJSON(http.StatusNotFound, "文件不存在")
		return
	}

	// 获取文件大小
	fileSize, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "无法获取文件大小")
		return
	}

	// 设置每个分片的大小（例如1MB）
	chunkSize := 1024 * 1024
	totalChunks := (fileSize + chunkSize - 1) / chunkSize

	// 创建用于存储下载分片的临时目录
	tempDir := filepath.Join("./tmp", "download_chunks")
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "创建临时目录失败")
		return
	}

	// 并发下载各个分片
	type chunkResult struct {
		index int
		err   error
	}
	resultChan := make(chan chunkResult, totalChunks)

	for i := 0; i < totalChunks; i++ {
		go func(index int) {
			start := index * chunkSize
			end := start + chunkSize - 1
			if end > fileSize-1 {
				end = fileSize - 1
			}

			// 设置Range请求头，只下载部分数据
			req, _ := http.NewRequest("GET", fileURL, nil)
			rangeHeader := fmt.Sprintf("bytes=%d-%d", start, end)
			req.Header.Set("Range", rangeHeader)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				resultChan <- chunkResult{index: index, err: err}
				return
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Println("关闭响应体失败:", err)
				}
			}(resp.Body)

			// 创建分片文件
			chunkFileName := filepath.Join(tempDir, strconv.Itoa(index))
			chunkFile, err := os.Create(chunkFileName)
			if err != nil {
				resultChan <- chunkResult{index: index, err: err}
				return
			}
			defer func(chunkFile *os.File) {
				err := chunkFile.Close()
				if err != nil {
					log.Println("关闭分片文件失败:", err)
				}
			}(chunkFile)

			// 写入分片数据
			_, err = io.Copy(chunkFile, resp.Body)
			resultChan <- chunkResult{index: index, err: err}
		}(i)
	}

	// 等待所有分片下载完成
	for i := 0; i < totalChunks; i++ {
		result := <-resultChan
		if result.err != nil {
			log.Printf("下载分片 %d 失败: %v\n", result.index, result.err)
			ctx.ErrorJSON(http.StatusInternalServerError, "下载失败")
			return
		}
	}

	// 合并所有分片
	finalFileName := filepath.Join("./downloads", filepath.Base(fileURL))
	finalFile, err := os.Create(finalFileName)
	if err != nil {
		ctx.ErrorJSON(http.StatusInternalServerError, "创建最终文件失败")
		return
	}
	defer func(finalFile *os.File) {
		err := finalFile.Close()
		if err != nil {
			log.Println("关闭最终文件失败:", err)
		}
	}(finalFile)

	for i := 0; i < totalChunks; i++ {
		chunkFileName := filepath.Join(tempDir, strconv.Itoa(i))
		chunkFile, err := os.Open(chunkFileName)
		if err != nil {
			ctx.ErrorJSON(http.StatusInternalServerError, "读取分片文件失败")
			return
		}

		_, err = io.Copy(finalFile, chunkFile)
		if err != nil {
			err := chunkFile.Close()
			if err != nil {
				return
			}
			ctx.ErrorJSON(http.StatusInternalServerError, "合并文件失败")
			return
		}
		err = chunkFile.Close()
		if err != nil {
			return
		}
		// 删除已合并的分片
		_ = os.Remove(chunkFileName)
	}

	// 清理临时目录
	_ = os.RemoveAll(tempDir)

	// 返回下载完成响应
	ctx.SuccessJSON("文件下载成功", map[string]string{"path": finalFileName})
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
