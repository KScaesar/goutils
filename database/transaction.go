package database

import (
	"context"
)

// TxFactory that lifecycle is equal to process scope, and is goroutine safe.
//
// 參數 ctx allow nil,
// 若 ctx 內部 有 tx 元件, 則使用 原本的 tx;
// 若 ctx 內部 無 tx 元件, 則產生 全新的 tx, 並將 tx 元件 assign 到 txCtx.
type TxFactory interface {
	CreateTx(ctx context.Context) Transaction
}

// Transaction lifecycle is equal to request scope, and not goroutine safe.
// You should use only one Transaction per goroutine
//
// 參數 fn 中的 txCtx,
// 保證 txCtx 一定有 tx 元件,
// 比如 *gorm.DB, mongo.Session.
type Transaction interface {
	// AutoComplete 執行 fn 之後, if fn return nil 自動執行 Commit, else return err 自動執行 Rollback
	AutoComplete(fn func(txCtx context.Context) error) error

	// ManualComplete 執行 fn 之後, 需要另外手動呼叫 Commit or Rollback
	ManualComplete(fn func(txCtx context.Context) error) (
		commit func() error,
		rollback func() error,
		fnErr error,
	)
}
