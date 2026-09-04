package https

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"samples/internal/components/logger"
	"samples/internal/utils/config"
	"strings"
	"time"
)

// initHttpClient 初始化HTTP客户端
func initHttpClient(cfg config.HttpConfig) {
	mu.Lock()
	defer mu.Unlock()
	client = &http.Client{
		Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second,
	}
}

// Get 发送GET请求
// requestUrl: 请求的URL
// header: 请求头
// resp: 响应体，传入一个指针类型的变量，用于接收响应数据
// 返回值: error，表示请求是否成功
func Get[T any](requestUrl string, header map[string]string, resp *T) (int, error) {
	return request(http.MethodGet, requestUrl, nil, header, resp)
}

// Post 发送POST请求
// requestUrl: 请求的URL
// body: 请求体
// header: 请求头
// resp: 响应体，传入一个指针类型的变量，用于接收响应数据
// 返回值: error，表示请求是否成功
func Post[T any](requestUrl string, body any, header map[string]string, resp *T) (int, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if header == nil {
		header = make(map[string]string)
	}
	header["Content-Type"] = "application/json"
	return request(http.MethodPost, requestUrl, strings.NewReader(string(data)), header, resp)
}

// PostForm 发送POST请求，使用application/x-www-form-urlencoded格式
// requestUrl: 请求的URL
// formData: 表单数据，传入一个map[string]string类型的变量，用于发送表单数据
// header: 请求头
// resp: 响应体，传入一个指针类型的变量，用于接收响应数据
// 返回值: error，表示请求是否成功
func PostForm[T any](requestUrl string, formData, header map[string]string, resp *T) (int, error) {
	form := url.Values{}
	for key, value := range formData {
		form.Set(key, value)
	}
	if header == nil {
		header = make(map[string]string)
	}
	header["Content-Type"] = "application/x-www-form-urlencoded"
	return request(http.MethodPost, requestUrl, strings.NewReader(form.Encode()), header, resp)
}

// PostFile 发送POST请求，使用multipart/form-data格式上传文件
// requestUrl: 请求的URL
// fileFieldName: 文件字段名
// filePath: 文件路径
// formData: 表单数据，传入一个map[string]string类型的变量，用于发送表单数据
// header: 请求头
// resp: 响应体，传入一个指针类型的变量，用于接收响应数据
// 返回值: error，表示请求是否成功
func PostFile[T any](requestUrl, fileFieldName, filePath string, formData, header map[string]string, resp *T) (int, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 文本字段
	for key, value := range formData {
		_ = writer.WriteField(key, value)
	}

	// 文件字段
	filePaths := strings.Split(filePath, "/")
	fileWriter, err := writer.CreateFormFile(fileFieldName, filePaths[len(filePaths)-1])
	if err != nil {
		return http.StatusInternalServerError, err
	}
	f, err := os.Open(filePath)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer f.Close()
	_, _ = io.Copy(fileWriter, f)

	_ = writer.Close()

	if header == nil {
		header = make(map[string]string)
	}
	header["Content-Type"] = writer.FormDataContentType()
	return request(http.MethodPost, requestUrl, &buf, header, resp)
}

// Put 发送PUT请求
// requestUrl: 请求的URL
// body: 请求体
// header: 请求头
// resp: 响应体，传入一个指针类型的变量，用于接收响应数据
// 返回值: error，表示请求是否成功
func Put[T any](requestUrl string, body any, header map[string]string, resp *T) (int, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if header == nil {
		header = make(map[string]string)
	}
	header["Content-Type"] = "application/json"
	return request(http.MethodPut, requestUrl, bytes.NewReader(data), header, resp)
}

// Delete 发送DELETE请求
// requestUrl: 请求的URL
// header: 请求头
// resp: 响应体，传入一个指针类型的变量，用于接收响应数据
// 返回值: error，表示请求是否成功
func Delete[T any](requestUrl string, header map[string]string, resp *T) (int, error) {
	return request(http.MethodDelete, requestUrl, nil, header, resp)
}

// request 发送HTTP请求
// method: 请求方法，如GET、POST、PUT、DELETE等
// requestUrl: 请求的URL
// body: 请求体，传入一个io.Reader类型的变量，用于发送请求数据
// header: 请求头
// resp: 响应体，传入一个指针类型的变量，用于接收响应数据
// 返回值: error，表示请求是否成功
func request[T any](method, requestUrl string, body io.Reader, header map[string]string, resp *T) (int, error) {
	req, err := http.NewRequest(method, requestUrl, body)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	if client == nil {
		return http.StatusInternalServerError, errors.New("HTTP客户端未初始化")
	}
	if body, ok := body.(*strings.Reader); ok {
		data := make([]byte, body.Size())
		_, _ = body.ReadAt(data, 0)
		logger.Debug("发送HTTP请求: %s %s, 请求体: %s\n", method, requestUrl, string(data))
	} else {
		logger.Debug("发送HTTP请求: %s %s\n", method, requestUrl)
	}
	httpResp, err := client.Do(req)
	if err != nil {
		return httpResp.StatusCode, err
	}
	defer httpResp.Body.Close()
	err = json.NewDecoder(httpResp.Body).Decode(resp)
	if err != nil {
		return httpResp.StatusCode, err
	}
	respData, _ := json.Marshal(resp)
	logger.Debug("接收HTTP响应: %s %s, 状态码: %d, 响应体: %s\n", method, requestUrl, httpResp.StatusCode, string(respData))
	return httpResp.StatusCode, nil
}
