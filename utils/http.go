package utils

import (
	"bytes"
	"encoding/json"
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

	rc, err := HttpProxyGet(url, nil, proxyURL)
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

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36"

// HttpGet获取指定的资源。如果是，则返回ErrNotFound
// 服务器以状态404响应。
func HttpProxyGet(rawurl string, header http.Header, proxyURL string) (io.ReadCloser, error) {
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

	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		err = errors.New("请求未发现")
	} else {
		err = errors.New("请求错误")
	}

	return nil, err
}

// http.NewRequest发送请求
//
// method:	请求方法：POST、GET、PUT、DELETE
// urlText:		请求地址
// contentType: 请求数据类型，首字母简写，如：axwfu
// params: 	请求提交的数据
// header:	请求头
func HttpRequest(method, urlText, contentType string, params, header map[string]string) (*http.Response, error) {
	if urlText == "" {
		panic(errors.New("url不能为空"))
	}
	method = strings.ToUpper(method)

	var req *http.Request
	var err error
	var body io.Reader
	if params != nil {
		if method == "POST" {
			switch contentType {
			case "axwfu": // application/x-www-form-urlencoded;
				data := make(url.Values)
				for k, v := range params {
					data[k] = []string{v}
				}
				body = strings.NewReader(data.Encode())
				contentType = "application/x-www-form-urlencoded; charset=utf-8"
			case "mfd": // multipart/form-data
				data := url.Values{}
				for k, v := range params {
					data.Set(k, v)
				}
				body = strings.NewReader(data.Encode())
				contentType = "multipart/form-data; charset=utf-8"
			case "tx": // text/xml
				data := url.Values{}
				for k, v := range params {
					data.Set(k, strings.ReplaceAll(v, " ", "+"))
				}
				body = strings.NewReader(data.Encode())
				contentType = "text/xml; charset=utf-8"
			default: // application/json
				jsonStr, err := json.Marshal(params)
				if err != nil {
					return nil, err
				}
				body = bytes.NewBuffer(jsonStr)
				contentType = "application/json; charset=utf-8"
			}
		} else {
			urlText = urlText + "?"
			for key, value := range params {
				urlText += key + "=" + value + "&"
			}
		}
		// url编码
		//urlText=urlText.QueryEscape(urlText)
	}
	req, err = http.NewRequest(method, urlText, body)

	if err != nil {
		return nil, err
	}
	if header != nil {
		for key, value := range header {
			req.Header.Add(key, value)
		}
	}
	if req.Header.Get("Content-Type") == "" && method == "POST" {
		req.Header.Set("Content-Type", contentType)
	}
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", UserAgent)
	}
	// dump出远程服务器返回的信息，调试请求
	//bd, err := httputil.DumpRequest(req, true)

	client := &http.Client{Timeout: 30 * time.Second}
	// 发起请求
	return client.Do(req)
}
