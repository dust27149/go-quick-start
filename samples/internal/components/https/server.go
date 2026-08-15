package https

import (
	"errors"
	"fmt"
	"net/http"
	"samples/internal/components/logger"
	"samples/internal/utils/config"
)

// registerRoutes 注册路由
func RegisterRoute(method, path string, handler http.HandlerFunc) {
	mux.HandleFunc(fmt.Sprintf("%s %s", method, path), handler)
}

// startServer 启动 HTTP 服务
func startServer(cfg config.HttpConfig) {
	mu.Lock()
	defer mu.Unlock()
	if server != nil {
		return
	}

	mux = http.NewServeMux()
	server = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("API 服务启动失败: %v", err)
		}
	}()
}
