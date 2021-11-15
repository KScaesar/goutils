package xLog

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

const CustomTimeFormat = "2006-01-02 15:04:05-07:00"

func init() {
	zerolog.ErrorStackMarshaler = errorStackMarshaler
	zerolog.ErrorMarshalFunc = errorMarshalFunc
	zerolog.TimestampFieldName = "timestamp"

	Init("debug", true)
}

// Init
// @param level: "debug", "info", "error", "panic"
func Init(level string, isDevEnv bool) {
	err := SetGlobalLevel(level)
	if err != nil {
		panic(err)
	}

	var w io.Writer
	if isDevEnv {
		zerolog.TimeFieldFormat = CustomTimeFormat
		w = NewConsoleWriter(os.Stdout)
	} else {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
		w = os.Stdout
	}
	SetDefaultLogger(NewLogger(w))
}

func NewConsoleWriter(w io.Writer) io.Writer {
	return &zerolog.ConsoleWriter{
		Out:        w,
		TimeFormat: CustomTimeFormat,
		FormatCaller: func(i interface{}) string {
			if i == nil { // 沒啟用 Caller 功能時, i == nil, 導致 i.(string) 發生錯誤
				return ""
			}
			return i.(string)
		},
	}
}
