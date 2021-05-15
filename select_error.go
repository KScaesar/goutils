package goutils

import (
	"context"

	"github.com/Min-Feng/goutils/errorY"
	"github.com/Min-Feng/goutils/logY"
)

// SelectError 從 ctx 取得 log 進行 abandon error 的 stack 紀錄
func SelectError(ctx context.Context, major error, abandon error) error {
	log := logY.FromCtx(ctx)

	abandonErr := errorY.WrapMessage(abandon, "[abandon error]")
	log.Err(abandonErr).Caller(1).Send()

	return errorY.WrapMessage(major, "[select error]")
}
