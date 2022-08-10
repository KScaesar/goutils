package xHttp

import (
	"github.com/KScaesar/goutils/errors"
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
		Code:    errors.Code(err),
		Message: err.Error(),
		Data:    struct{}{},
	}
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}
