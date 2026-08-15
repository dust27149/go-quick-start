package https

import (
	"net/http"
	"samples/internal/utils/config"
	"sync"
)

var mu sync.Mutex       // 互斥锁，确保并发访问安全
var client *http.Client // 定义一个全局的HTTP客户端，统一设置超时时间等参数
var server *http.Server // 定义一个全局的HTTP服务实例
var mux *http.ServeMux  // 定义一个全局的HTTP请求多路复用器

// Init 初始化HTTP客户端
func Init(cfg config.HttpConfig) {
	// 初始化HTTP客户端
	initHttpClient(cfg)
	// 启动 HTTP 服务
	startServer(cfg)
	// 注册API路由
	registerBasicApi()
}

// DeInit 关闭HTTP客户端，释放资源
func DeInit() {
	mu.Lock()
	defer mu.Unlock()
	if client != nil {
		client.CloseIdleConnections()
		client = nil
	}
	if server != nil {
		server.Close()
		server = nil
	}
}
