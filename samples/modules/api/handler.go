package api

import (
	"fmt"
	"net/http"
	"samples/internal/components/https"
)

// testHandler 处理 GET /test 路由的请求
func testHandler(w http.ResponseWriter, r *http.Request) {
	for key, values := range https.ReadRequestParams(r) {
		fmt.Printf("DEBUG 查询参数: %s=%s\n", key, values)
	}
	https.WriteResponse(w, 200, https.SuccessResponse())
}
