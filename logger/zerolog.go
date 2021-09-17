package logger

import (
	"github.com/rs/zerolog"
)

// 為了避免其他使用者 import zerolog
type (
	ZeroLogger     = zerolog.Logger
	ZeroLogContext = zerolog.Context
	UpdateContext  func(c ZeroLogContext) ZeroLogContext
)
