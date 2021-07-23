package database

import (
	"context"
)

// UowFactory that lifecycle is equal to process scope
type UowFactory interface {
	CreateUow() (Uow, error)
}

// Uow is unit of work, 類似原子性的概念, ref: https://martinfowler.com/eaaCatalog/unitOfWork.html .
// and lifecycle is equal to request scope
//
// 參數 txCtx 表示內部含有 執行 transaction 的物件, 比如 *gorm.DB, mongo.Session
// 參數 ctx allow nil
//
// AutoStart 執行 txFn 之後, if txFn return nil 自動執行 Commit, else return err 自動執行 Rollback
//
// ManualStart 執行 txFn 之後, 需要另外手動呼叫 Commit or Rollback
type Uow interface {
	AutoStart(ctx context.Context, txFn func(txCtx context.Context) error) error

	ManualStart(ctx context.Context, txFn func(txCtx context.Context) error) error
	Commit() error
	Rollback() error
}
