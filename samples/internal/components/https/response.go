package https

// Response 接收不包含data或list字段的响应
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// DataResponse 接收包含data字段的响应
type DataResponse[T any] struct {
	Response
	Data T `json:"data"`
}

// ListResponse 接收包含list字段的响应
type ListResponse[T any] struct {
	Response
	List []T `json:"list"`
}

// ErrorResponse 返回一个错误响应对象
func ErrorResponse(code int, msg string) Response {
	return Response{
		Code:    code,
		Message: msg,
		Success: false,
	}
}

// SuccessResponse 返回一个成功响应对象
func SuccessResponse() Response {
	return Response{
		Code:    0,
		Message: "请求成功",
		Success: true,
	}
}

// ErrorDataResponse 返回一个包含data字段的错误响应对象
func ErrorDataResponse[T any](code int, msg string) DataResponse[T] {
	return DataResponse[T]{
		Response: Response{
			Code:    code,
			Message: msg,
			Success: false,
		},
	}
}

// SuccessDataResponse 返回一个包含data字段的成功响应对象
func SuccessDataResponse[T any](data T) DataResponse[T] {
	return DataResponse[T]{
		Response: Response{
			Code:    CODE_SUCCESS,
			Message: "请求成功",
			Success: true,
		},
		Data: data,
	}
}

// ErrorListResponse 返回一个包含list字段的错误响应对象
func ErrorListResponse[T any](code int, msg string) ListResponse[T] {
	return ListResponse[T]{
		Response: Response{
			Code:    code,
			Message: msg,
			Success: false,
		},
	}
}

// SuccessListResponse 返回一个包含list字段的成功响应对象
func SuccessListResponse[T any](list []T) ListResponse[T] {
	return ListResponse[T]{
		Response: Response{
			Code:    CODE_SUCCESS,
			Message: "请求成功",
			Success: true,
		},
		List: list,
	}
}
