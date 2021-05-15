package logY

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

	DefaultMode()
}

func DefaultMode() {
	Init(Config{Show: "human", Level: "info"})
}

func FixBugMode() {
	SetGlobalLevel("debug")
}

// TestingMode 避免執行 go test, 出現 log 訊息
func TestingMode() {
	SetGlobalLevel("panic")
}

type Config struct {
	Show  string // "json", "human"
	Level string // "debug", "info", "error", "panic"
}

func Init(cfg Config) {
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
		_default = New(NewConsoleWriter(os.Stdout))

	default:
		// 這裡簡化了參數, show 包含了 json 格式 及 io.Writer
		// 為了讓之後 io.Writer, 可以替換
		// 所以略過 default 設定
	}
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
