package database

import (
	"context"
)

// TxFactory that lifecycle is equal to process scope
type TxFactory interface {
	CreateTx() (Transaction, error)
}

// Transaction lifecycle is equal to request scope
//
// 參數 txCtx 表示內部含有 執行 transaction 的物件, 比如 *gorm.DB, mongo.Session
// 參數 ctx allow nil
type Transaction interface {
	// AutoStart 執行 fn 之後, if fn return nil 自動執行 Commit, else return err 自動執行 Rollback
	AutoStart(ctx context.Context, fn func(txCtx context.Context) error) error

	// ManualStart 執行 fn 之後, 需要另外手動呼叫 Commit or Rollback
	ManualStart(
		ctx context.Context,
		fn func(txCtx context.Context) error,
	) (
		commit func() error,
		rollback func() error,
		fnErr error,
	)
}
