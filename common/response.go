package common

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func Success(message string) *Response {
	return &Response{
		Code:    200,
		Message: message,
		Type:    "success",
		Success: true,
		Data:    nil,
	}
}

func SuccessWithData(message string, data interface{}) *Response {
	return &Response{
		Code:    200,
		Message: message,
		Type:    "success",
		Success: true,
		Data:    data,
	}
}

func Warning(message string) *Response {
	return &Response{
		Code:    200,
		Message: message,
		Type:    "warning",
		Success: false,
		Data:    nil,
	}
}

func Error(message string) *Response {
	return &Response{
		Code:    200,
		Message: message,
		Type:    "error",
		Success: false,
		Data:    nil,
	}
}

func Fatal(message string) *Response {
	return &Response{
		Code:    500,
		Message: message,
		Type:    "error",
		Success: false,
		Data:    nil,
	}
}
