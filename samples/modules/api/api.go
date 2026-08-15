package api

import "samples/internal/components/https"

// Register 注册路由
func Register() {
	https.RegisterRoute("GET", "/test", testHandler)
}
