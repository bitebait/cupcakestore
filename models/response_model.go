package models

type Response struct {
	Error   bool        `json:"error"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func NewResponse(error bool, data interface{}, message string) *Response {
	return &Response{
		Error:   error,
		Data:    data,
		Message: message,
	}
}
