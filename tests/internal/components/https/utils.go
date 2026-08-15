package https

import (
	"encoding/json"
	"net/http"
)

// ReadHeader 从请求头中提取指定的值
func ReadHeader(r *http.Request, key string) string {
	return r.Header.Get(key)
}

// ReadRequestParam 从 URL 查询参数中提取指定的参数值
func ReadRequestParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// ReadRequestParams 从 URL 查询参数中提取所有参数值
func ReadRequestParams(r *http.Request) map[string][]string {
	return r.URL.Query()
}

// ReadFormValue 从表单数据中提取指定的值
func ReadFormValue(r *http.Request, key string) string {
	return r.FormValue(key)
}

// ReadRequestBody 从请求体中读取 JSON 数据并解析为 map
func ReadRequestBody(r *http.Request) (map[string]interface{}, error) {
	var requestBody map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return nil, err
	}
	return requestBody, nil
}

// WriteResponse 返回响应给客户端
func WriteResponse(w http.ResponseWriter, code int, resp Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "服务端错误", http.StatusInternalServerError)
	}
}
