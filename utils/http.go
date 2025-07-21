package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	ContentTypeAXWFU = "application/x-www-form-urlencoded"
	ContentTypeMFD   = "multipart/form-data"
	ContentTypeTX    = "text/xml"
	ContentTypeJson  = "application/json"
)

// UserAgent
// https://useragentstring.com
// https://www.useragents.me
// https://www.whatismybrowser.com
// https://explore.whatismybrowser.com/useragents/explore
// https://www.user-agents.org/allagents.xml
// https://techpatterns.com/downloads/firefox/useragentswitcher.xml
var UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36"

// DownFile 下载文件
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

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	_, err = io.Copy(f, rc)

	return uploadDir + newName, err
}

// HttpProxyGet HttpGet获取指定的资源。如果是，则返回ErrNotFound
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode == 200 {
		return resp.Body, nil
	}
	if resp.StatusCode == 404 {
		err = errors.New("请求未发现")
	} else {
		err = errors.New("请求错误")
	}

	return nil, err
}

// contextCancelingBody 是一个包装器，它包装了 http.Response.Body。
// 它的目的是在 Body 被关闭或被完全读取后，自动调用 context 的 cancel 函数。
//
// 这里也可以处理Body在读取后io.ReadCloser，底层的网络数据流已经被读取到了末尾（EOF），读取一次后即被消费的问题
// 创建一个新的 Body
// bodyBytes, err := io.ReadAll(resp.Body)
// resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
// reader := bytes.NewReader(bodyBytes)
type contextCancelingBody struct {
	io.ReadCloser                    // 嵌入原始的 Body，继承其 Read 和 Close 方法
	cancel        context.CancelFunc // 持有需要被调用的 cancel 函数
}

// Read 重写 Read 方法
// 当从原始 Body 读取数据时，我们会检查是否读到了末尾 (io.EOF)。
// 如果读到了末尾，说明 Body 已被耗尽，我们就可以安全地取消 context 了。
func (b *contextCancelingBody) Read(p []byte) (n int, err error) {
	n, err = b.ReadCloser.Read(p)
	if err == io.EOF {
		// 读取到流的末尾，意味着工作完成，调用cancel
		b.cancel()
	}
	return n, err
}

// Close 重写 Close 方法
// 调用者有责任关闭 Body。当他们调用 Close 时，也意味着他们不再需要这个 Body 了，
// 我们可以安全地取消 context，并调用原始 Body 的 Close 方法。
func (b *contextCancelingBody) Close() error {
	// 调用 cancel 是幂等的（调用多次也没关系），所以在这里调用是安全的
	b.cancel()
	return b.ReadCloser.Close()
}

// HttpRequest 使用http.DefaultClient复用TCP连接池极大地减少了网络握手的开销，http.NewRequest发送请求
//
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
					if err := writer.WriteField(k, v); err != nil {
						return nil, err
					}
				}
				if err := writer.Close(); err != nil {
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
	// 为请求添加超时
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	req, err := http.NewRequestWithContext(ctx, method, urlText, body)
	//req, err := http.NewRequest(method, urlText, body)

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

	// http.DefaultClient 已经内置了一个 http.Transport，它会维护一个可复用的TCP连接池 (Connection Pooling / Keep-Alive)
	// 频繁创建新的 http.Client 会失去这个优势
	// 发起请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	// 创建一个包装器，它包装了原始的 Body，并添加了 cancel 函数
	resp.Body = &contextCancelingBody{
		ReadCloser: resp.Body,
		cancel:     cancel,
	}
	return resp, err
}

// HttpReadBody 请求并读取返回内容
func HttpReadBody(method, urlText, contentType string, params, header map[string]string) ([]byte, error) {
	res, err := HttpRequest(method, urlText, contentType, params, header)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)
	result, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// HttpReadBodyString 请求并读取返回内容为字符串
func HttpReadBodyString(method, urlText, contentType string, params, header map[string]string) (string, error) {
	res, err := HttpRequest(method, urlText, contentType, params, header)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)
	result, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// HttpReadBodyJsonObject 请求并读取返回内容为json对象
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

// HttpReadBodyJsonMap 请求并读取返回内容为json map
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

// HttpReadBodyJsonMapArray 请求并读取返回内容为json Map数组
func HttpReadBodyJsonMapArray(method, urlText, contentType string, params,
	header map[string]string) ([]map[string]interface{}, error) {
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

// HttpReadBodyJsonArray 请求并读取返回内容为json数组
func HttpReadBodyJsonArray(method, urlText, contentType string, params, header map[string]string) ([]interface{}, error) {
	res, err := HttpReadBody(method, urlText, contentType, params, header)
	if err != nil {
		return nil, err
	}
	var data []interface{}
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

func (hc *HttpClient) HttpReadBodyJsonMapArray() ([]map[string]interface{}, error) {
	return HttpReadBodyJsonMapArray(hc.Method, hc.UrlText, hc.ContentType, hc.Params, hc.Header)
}

func (hc *HttpClient) ReadBodyJsonArray() ([]interface{}, error) {
	return HttpReadBodyJsonArray(hc.Method, hc.UrlText, hc.ContentType, hc.Params, hc.Header)
}

// fetchURL 在给定的上下文环境中通过HTTP GET请求获取指定URL的响应内容。
// 该函数接受一个上下文对象、一个HTTP客户端对象和一个URL字符串作为参数。
// 它返回URL响应的字节切片和可能出现的错误。
func fetchURL(ctx context.Context, client *http.Client, url string) ([]byte, error) {
	// 创建一个HTTP GET请求，并将上下文对象嵌入其中。
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 使用提供的HTTP客户端对象执行请求。
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// 确保在函数返回前关闭响应体。
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	// 读取并返回响应体的全部内容。
	return io.ReadAll(resp.Body)
}

// FetchAll 并发地获取一组URL的内容，并以映射表的形式返回结果。
// 参数 urls 是待获取的URL列表。
// 返回值是一个映射表，键为URL，值为该URL对应的内容。
func FetchAll(urls []string) map[string][]byte {
	// 创建一个专门用于 HTTP/2 的 transport
	/*h2Transport := &http2.Transport{
		AllowHTTP: true, // 允许非加密的 h2c
		DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
			return tls.Dial(network, addr, cfg)
		},
	}*/
	// 创建一个可复用的、高度优化的HTTP客户端，配置连接池和请求超时。
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   5 * time.Second,  // 建立TCP连接的超时时间
				KeepAlive: 30 * time.Second, // TCP Keep-Alive 间隔
			}).DialContext,
			MaxIdleConns:        100,              // 最大空闲连接数
			MaxIdleConnsPerHost: 100,              // 对每个主机的最大空闲连接数
			IdleConnTimeout:     90 * time.Second, // 空闲连接在关闭前保持的时间
			TLSHandshakeTimeout: 10 * time.Second, // TLS 握手超时
		},
		Timeout: 20 * time.Second, // 整个请求的超时时间，包括连接、重定向、读取响应体
		/*CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},*/
	}

	// 创建一个带有超时的上下文，用于控制请求的最长执行时间。
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	results := make(map[string][]byte)
	mu := &sync.Mutex{}

	// 遍历URL列表，对每个URL发起并发请求。
	for _, url_ := range urls {
		wg.Add(1) // WaitGroup 计数器+1
		// 使用闭包捕获当前URL，避免循环变量复用问题。
		go func(u string) {
			defer wg.Done() // goroutine 结束时，计数器-1
			// 使用上下文和客户端获取URL内容。
			data, err := fetchURL(ctx, client, u)
			if err != nil {
				log.Printf("Error fetching %s: %v\n", u, err)
				return
			}
			// 使用互斥锁确保并发安全，然后将结果存入映射表中。
			mu.Lock()
			results[u] = data
			mu.Unlock()
		}(url_)
	}

	// 等待所有并发请求完成。
	wg.Wait()
	// 返回所有URL的内容。
	return results
}
