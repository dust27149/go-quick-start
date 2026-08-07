package client

import (
	"net/http"
	"time"
)

// NewHTTPClient 创建一个带统一超时时间的 HTTP 客户端。
// 统一从这里创建客户端，后面如果要加重试、Header 或代理配置会更容易扩展。
func NewHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout}
}
