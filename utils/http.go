package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
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

var UserAgent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/29.0.1541.0 Safari/537.36"

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

	resp.Body.Close()
	if resp.StatusCode == 404 {
		err = errors.New("请求未发现")
	} else {
		err = errors.New("请求错误")
	}

	return nil, err
}

// http.Client发送请求
// method:	请求方法：POST、GET、PUT、DELETE
// url:		请求地址
// params:	请求参数
func HttpClient(method, url string, params map[string]string) string {
	method = strings.ToUpper(method)

	client := http.Client{Timeout: 5 * time.Second}

	var resp *http.Response
	var err error
	if method == "GET" {
		param := "?"
		for key, value := range params {
			param += key + "=" + value + "&"
		}
		param = param[0 : len(param)-1]
		resp, err = client.Get(url + param)
	} else if method == "POST" {
		jsonStr, _ := json.Marshal(params)
		resp, err = client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	}

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	return result.String()
}

// http.NewRequest发送请求
// method:	请求方法：POST、GET、PUT、DELETE
// url:		请求地址
// params: 	请求提交的数据
// header:	请求体格式，如：application/json
func HttpRequest(method, url string, params map[string]string, header map[string]string) string {
	var req *http.Request
	var err error
	if params != nil {
		jsonStr, _ := json.Marshal(params)
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	} else {
		param := "?"
		for key, value := range params {
			param += key + "=" + value + "&"
		}
		param = param[0 : len(param)-1]
		req, err = http.NewRequest(method, url+param, nil)
	}

	if header != nil {
		for key, value := range header {
			req.Header.Add(key, value)
		}
	}

	if header == nil || req.Header.Get("content-type") == "" {
		req.Header.Add("content-type", "application/json")
	}
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, error := client.Do(req)
	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return string(result)
}
