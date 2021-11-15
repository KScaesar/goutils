package xLog

import (
	"github.com/rs/zerolog"

	"github.com/Min-Feng/goutils/errors"
)

func errorStackMarshaler(err error) interface{} {
	return errors.Stacks(err)
}

func errorMarshalFunc(err error) interface{} {
	return &errorObject{
		msg:  err.Error(),
		code: errors.Code(err),
	}
}

type errorObject struct {
	msg  string
	code int
}

func (obj *errorObject) MarshalZerologObject(e *zerolog.Event) {
	e.Str("err_msg", obj.msg).Int("err_code", obj.code)
}
