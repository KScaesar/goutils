package logger

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
	e.Str("msg", o.msg).Int("code", o.code)
}
