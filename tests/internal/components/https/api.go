package https

import (
	"fmt"
	"net/http"
)

// registerBasicApi 注册API路由
func registerBasicApi() {
	RegisterRoute("GET", "/health", getHealthHandler)
	RegisterRoute("POST", "/health", postHealthHandler)
}

// getHealthHandler 处理 GET /health 路由的请求
func getHealthHandler(w http.ResponseWriter, r *http.Request) {
	for key, values := range ReadRequestParams(r) {
		fmt.Printf("DEBUG 查询参数: %s=%s\n", key, values)
	}
	WriteResponse(w, 200, SuccessResponse())
}

// postHealthHandler 处理 POST /health 路由的请求
func postHealthHandler(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ReadRequestBody(r)
	if err != nil {
		WriteResponse(w, http.StatusBadRequest, ErrorResponse(1, "请求参数解析失败"))
		return
	}
	fmt.Printf("DEBUG 收到请求: %+v\n", requestBody)
	WriteResponse(w, 200, SuccessResponse())
}
