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

// errorResponse 返回一个错误响应对象
func ErrorResponse(code int, msg string) Response {
	return Response{
		Code:    code,
		Message: msg,
		Success: false,
	}
}

// successResponse 返回一个成功响应对象
func SuccessResponse() Response {
	return Response{
		Code:    200,
		Message: "请求成功",
		Success: true,
	}
}

// errorDataResponse 返回一个包含data字段的错误响应对象
func ErrorDataResponse[T any](code int, msg string) DataResponse[T] {
	return DataResponse[T]{
		Response: Response{
			Code:    code,
			Message: msg,
			Success: false,
		},
	}
}

// successDataResponse 返回一个包含data字段的成功响应对象
func SuccessDataResponse[T any](data T) DataResponse[T] {
	return DataResponse[T]{
		Response: Response{
			Code:    200,
			Message: "请求成功",
			Success: true,
		},
		Data: data,
	}
}

// errorListResponse 返回一个包含list字段的错误响应对象
func ErrorListResponse[T any](code int, msg string) ListResponse[T] {
	return ListResponse[T]{
		Response: Response{
			Code:    code,
			Message: msg,
			Success: false,
		},
	}
}

// successListResponse 返回一个包含list字段的成功响应对象
func SuccessListResponse[T any](list []T) ListResponse[T] {
	return ListResponse[T]{
		Response: Response{
			Code:    200,
			Message: "请求成功",
			Success: true,
		},
		List: list,
	}
}
