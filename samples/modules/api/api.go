package api

import (
	"net/http"
	"samples/internal/components/https"
)

// Register 注册路由
func Register() {
	https.RegisterRoute(http.MethodGet, "/test", testHandler)
}
