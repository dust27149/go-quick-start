package requests

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
	"tests/internal/utils/config"
	"time"
)

var client *http.Client // 定义一个全局的HTTP客户端，统一设置超时时间等参数
var mu sync.Mutex       // 互斥锁，确保对 httpClient 的并发访问安全

// Init 初始化HTTP客户端
func Init(cfg config.HttpConfig) {
	mu.Lock()
	defer mu.Unlock()

	client = &http.Client{
		Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second,
	}
}

// DeInit 关闭HTTP客户端，释放资源
func DeInit() {
	mu.Lock()
	defer mu.Unlock()
	if client != nil {
		client.CloseIdleConnections()
		client = nil
	}
}

// Get 发送GET请求
// url: 请求的URL
// header: 请求头
// resp: 响应体，传入一个指针类型的变量，用于接收响应数据
// 返回值: error，表示请求是否成功
func Get[T any](url string, header map[string]string, resp *T) error {
	return request("GET", url, nil, header, resp)
}

// Post 发送POST请求
// url: 请求的URL
// body: 请求体，传入一个io.Reader类型的变量，用于发送请求数据
// header: 请求头
// resp: 响应体，传入一个指针类型的变量，用于接收响应数据
// 返回值: error，表示请求是否成功
func Post[T any](url string, body io.Reader, header map[string]string, resp *T) error {
	return request("POST", url, body, header, resp)
}

// Put 发送PUT请求
// url: 请求的URL
// body: 请求体，传入一个io.Reader类型的变量，用于发送请求数据
// header: 请求头
// resp: 响应体，传入一个指针类型的变量，用于接收响应数据
// 返回值: error，表示请求是否成功
func Put[T any](url string, body io.Reader, header map[string]string, resp *T) error {
	return request("PUT", url, body, header, resp)
}

// Delete 发送DELETE请求
// url: 请求的URL
// header: 请求头
// resp: 响应体，传入一个指针类型的变量，用于接收响应数据
// 返回值: error，表示请求是否成功
func Delete[T any](url string, header map[string]string, resp *T) error {
	return request("DELETE", url, nil, header, resp)
}

// request 发送HTTP请求
// method: 请求方法，如GET、POST、PUT、DELETE等
// url: 请求的URL
// body: 请求体，传入一个io.Reader类型的变量，用于发送请求数据
// header: 请求头
// resp: 响应体，传入一个指针类型的变量，用于接收响应数据
// 返回值: error，表示请求是否成功
func request[T any](method, url string, body io.Reader, header map[string]string, resp *T) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	if client == nil {
		return errors.New("HTTP客户端未初始化")
	}
	httpResp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()
	return json.NewDecoder(httpResp.Body).Decode(resp)
}
