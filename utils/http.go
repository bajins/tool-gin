package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	ContentTypeAXWFU = "application/x-www-form-urlencoded"
	ContentTypeMFD   = "multipart/form-data"
	ContentTypeTX    = "text/xml"
	ContentTypeJson  = "application/json"
	UserAgent        = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36"
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
// Content-Type只会存在于POST、PATCH、PUT等请求有请求数据实体时指定数据类型和数据字符集编码，
// 而GET、DELETE、HEAD、OPTIONS、TRACE等请求没有请求数据实体
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
	if method != http.MethodGet && method != http.MethodDelete && method != http.MethodHead &&
		method != http.MethodOptions && method != http.MethodTrace && method != http.MethodPost &&
		method != http.MethodPatch && method != http.MethodPut {
		return nil, errors.New("method不正确")
	}
	var req *http.Request
	var err error
	var body io.Reader
	if params != nil {
		if method == http.MethodPost || method == http.MethodPatch || method == http.MethodPut {
			switch contentType {
			case ContentTypeAXWFU: // application/x-www-form-urlencoded;
				data := make(url.Values)
				//data := url.Values{}
				for k, v := range params {
					data[k] = []string{v}
					//data.Set(k, v)
				}
				body = strings.NewReader(data.Encode())
				contentType = "application/x-www-form-urlencoded; charset=utf-8"
			case ContentTypeMFD: // multipart/form-data
				bodyBuf := &bytes.Buffer{}
				writer := multipart.NewWriter(bodyBuf)
				for k, v := range params {
					if err = writer.WriteField(k, v); err != nil {
						return nil, err
					}
				}
				if err = writer.Close(); err != nil {
					return nil, err
				}
				body = bodyBuf
				contentType = writer.FormDataContentType()
			case ContentTypeTX: // text/xml
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
	if req.Header.Get("Content-Type") == "" && (method == http.MethodPost ||
		method == http.MethodPatch || method == http.MethodPut) {
		req.Header.Set("Content-Type", contentType)
	}
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", UserAgent)
	}
	// dump出远程服务器返回的信息，调试请求
	//bd, err := httputil.DumpRequest(req, true)

	client := &http.Client{Timeout: 3 * time.Minute}
	// 发起请求
	return client.Do(req)
}

// 请求并读取返回内容
func HttpReadBody(method, urlText, contentType string, params, header map[string]string) ([]byte, error) {
	res, err := HttpRequest(method, urlText, contentType, params, header)
	if err != nil {
		return nil, err
	}
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 请求并读取返回内容为字符串
func HttpReadBodyString(method, urlText, contentType string, params, header map[string]string) (string, error) {
	res, err := HttpRequest(method, urlText, contentType, params, header)
	if err != nil {
		return "", err
	}
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// 请求并读取返回内容为json对象
func HttpReadBodyJsonObject(method, urlText, contentType string, params, header map[string]string, obj *interface{}) error {
	res, err := HttpReadBody(method, urlText, contentType, params, header)
	if err != nil {
		return err
	}
	err = json.Unmarshal(res, obj)
	if err != nil {
		return err
	}
	return nil
}

// 请求并读取返回内容为json map
func HttpReadBodyJsonMap(method, urlText, contentType string, params, header map[string]string) (map[string]interface{}, error) {
	res, err := HttpReadBody(method, urlText, contentType, params, header)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// 请求并读取返回内容为json数组
func HttpReadBodyJsonArray(method, urlText, contentType string, params, header map[string]string) ([]map[string]interface{}, error) {
	res, err := HttpReadBody(method, urlText, contentType, params, header)
	if err != nil {
		return nil, err
	}
	var data []map[string]interface{}
	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type HttpClient struct {
	Method      string
	UrlText     string
	ContentType string
	Params      map[string]string
	Header      map[string]string
}

func (hc *HttpClient) HttpRequest() (*http.Response, error) {
	return HttpRequest(hc.Method, hc.UrlText, hc.ContentType, hc.Params, hc.Header)
}

func (hc *HttpClient) ReadBody() ([]byte, error) {
	return HttpReadBody(hc.Method, hc.UrlText, hc.ContentType, hc.Params, hc.Header)
}

func (hc *HttpClient) ReadBodyString() (string, error) {
	return HttpReadBodyString(hc.Method, hc.UrlText, hc.ContentType, hc.Params, hc.Header)
}

func (hc *HttpClient) HttpReadBodyJsonObject(obj *interface{}) error {
	return HttpReadBodyJsonObject(hc.Method, hc.UrlText, hc.ContentType, hc.Params, hc.Header, obj)
}

func (hc *HttpClient) HttpReadBodyJsonMap() (map[string]interface{}, error) {
	return HttpReadBodyJsonMap(hc.Method, hc.UrlText, hc.ContentType, hc.Params, hc.Header)
}

func (hc *HttpClient) ReadBodyJsonArray() ([]map[string]interface{}, error) {
	return HttpReadBodyJsonArray(hc.Method, hc.UrlText, hc.ContentType, hc.Params, hc.Header)
}
