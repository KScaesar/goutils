package xLog

import (
	"github.com/rs/zerolog"

	"github.com/KScaesar/goutils/errors"
)

func errorStackMarshaler(err error) interface{} {
	if IsDebugLevel() {
		return errors.Stacks(err)
	}
	return nil
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
