package utils

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// 下载文件
func DownFile(url, upPreDir, upDir string, proxyURL string) (string, error) {
	fileType := url[strings.LastIndex(url, "."):]
	if fileType != ".jpeg" && fileType != ".png" && fileType != ".jpg" {
		fileType = ".jpeg"
	}
	newName := strconv.FormatInt(time.Now().UnixNano(), 10) + fileType
	uploadDir := upDir + time.Now().Format("2006/01/02") + "/"

	err := os.MkdirAll(upPreDir+uploadDir, os.ModePerm) //创建目录
	if err != nil {
		return "", err
	}

	rc, err := HttpGet(url, nil, proxyURL)
	if err != nil {
		return "", err
	}

	f, err := os.Create(upPreDir + uploadDir + newName)
	if err != nil {
		return "", err
	}

	defer f.Close()
	_, err = io.Copy(f, rc)

	return uploadDir + newName, err
}

var UserAgent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/29.0.1541.0 Safari/537.36"

// HttpGet获取指定的资源。如果是，则返回ErrNotFound
// 服务器以状态404响应。
func HttpGet(rawurl string, header http.Header, proxyURL string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", rawurl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Proxy-Switch-Ip", "yes")
	req.Header.Set("User-Agent", UserAgent)
	for k, vs := range header {
		req.Header[k] = vs
	}

	// 设置请求超时时间
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	// 设置代理
	if proxyURL != "" {
		parsedProxyUrl, err := url.Parse(proxyURL)
		if err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(parsedProxyUrl),
			}
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		return resp.Body, nil
	}

	resp.Body.Close()
	if resp.StatusCode == 404 { //
		err = errors.New("请求未发现")
	} else {
		err = errors.New("请求错误")
	}

	return nil, err
}
