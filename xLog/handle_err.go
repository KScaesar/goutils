package xLog

import (
	"github.com/rs/zerolog"

	"github.com/Min-Feng/goutils/errors"
)

func errorStackMarshaler(err error) interface{} {
	return errors.Stacks(err)
}

func errorMarshalFunc(err error) interface{} {
	return &errorLogObject{
		msg:  err.Error(),
		code: errors.Code(err),
	}
}

type errorLogObject struct {
	msg  string
	code int
}

func (o *errorLogObject) MarshalZerologObject(e *zerolog.Event) {
	e.Str("err_msg", o.msg).Int("err_code", o.code)
}
