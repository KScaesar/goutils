package httpY

import (
	"github.com/Min-Feng/goutils/errorY"
)

func NewNormalResponse(data interface{}) *Response {
	if data == nil {
		data = struct{}{}
	}
	return &Response{
		Code:    0,
		Message: "ok",
		Data:    data,
	}
}

func NewErrorResponse(err error) *Response {
	return &Response{
		Code:    errorY.Code(err),
		Message: errorY.SimpleInfo(err),
		Data:    struct{}{},
	}
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}
