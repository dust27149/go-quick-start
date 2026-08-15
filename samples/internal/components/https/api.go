package https

import (
	"net/http"
	"strings"
)

// registerBasicApi 注册API路由
func registerBasicApi() {
	RegisterRoute(http.MethodGet, "/health", getHealthHandler)
	RegisterRoute(http.MethodPost, "/health", postHealthHandler)
}

// getHealthHandler 处理 GET /health 路由的请求
func getHealthHandler(w http.ResponseWriter, r *http.Request) {
	params := make(map[string]string)
	for key, values := range ReadRequestParams(r) {
		params[key] = strings.Join(values, ",")
	}
	WriteResponse(w, 200, SuccessDataResponse(params))
}

// postHealthHandler 处理 POST /health 路由的请求
func postHealthHandler(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ReadRequestBody(r)
	if err != nil {
		WriteResponse(w, http.StatusBadRequest, ErrorResponse(1, "请求参数解析失败"))
		return
	}
	WriteResponse(w, 200, SuccessDataResponse(requestBody))
}
