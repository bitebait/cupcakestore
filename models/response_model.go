package models

type Response struct {
	Error   bool        `json:"error"`
	Object  interface{} `json:"object"`
	Message string      `json:"message"`
}

func NewResponse(error bool, object interface{}, message string) *Response {
	return &Response{
		Error:   error,
		Object:  object,
		Message: message,
	}
}
