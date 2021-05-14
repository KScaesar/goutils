package logY

import (
	"io"
	"os"

	"github.com/Min-Feng/goutils/errorY"

	"github.com/rs/zerolog"
)

const CustomTimeFormat = "2006-01-02 15:04:05-07:00"

func init() {
	DefaultMode()
}

func DefaultMode() {
	Init(Config{Show: "human", Level: "info"})
}

func FixBugMode() {
	SetGlobalLevel("debug")
}

func TestingMode() {
	SetGlobalLevel("panic")
}

type Config struct {
	Show  string // "json", "human"
	Level string // "debug", "info", "error", "panic"
}

func Init(cfg Config) {
	zerolog.ErrorStackMarshaler = errMarshalStack
	zerolog.TimestampFieldName = "timestamp"

	err := SetGlobalLevel(cfg.Level)
	if err != nil {
		panic(err)
	}

	switch cfg.Show {
	case "json":
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
		_default = New(os.Stdout)

	case "human":
		zerolog.TimeFieldFormat = CustomTimeFormat
		_default = New(consoleWriter())

	default:
		// 這裡簡化了參數, show 包含了 json 格式 及 io.Writer
		// 為了讓之後 io.Writer, 可以替換
		// 所以略過 default 設定
	}
}

func errMarshalStack(err error) interface{} {
	return errorY.Stacks(err)
}

func consoleWriter() io.Writer {
	writer := &zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: CustomTimeFormat,
		FormatCaller: func(i interface{}) string {
			if i == nil { // 沒啟用 Caller 功能時, i == nil, 導致 i.(string) 發生錯誤
				return ""
			}
			return i.(string)
		},
	}
	return writer
}
