package goutils

import (
	"github.com/Min-Feng/goutils/errorY"
	"github.com/Min-Feng/goutils/logY"
)

func SelectError(log *logY.WrapperLogger, major error, abandon error) error {
	if log == nil {
		freshLog := logY.Logger()
		log = &freshLog
	}

	abandonErr := errorY.WrapMessage(abandon, "[abandon error]")
	log.Err(abandonErr).Caller(1).Send()

	return errorY.WrapMessage(major, "[select error]")
}
