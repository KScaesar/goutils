package database

import (
	"context"

	"gorm.io/gorm"

	"github.com/Min-Feng/goutils/errorY"
)

func NewGormTxFactory(db *WrapperGorm) TxFactory {
	return &gormTxFactory{wDB: db}
}

type gormTxFactory struct {
	wDB *WrapperGorm
}

func (f *gormTxFactory) CreateTx() (Transaction, error) {
	return &gormTxAdapter{wrapperDB: f.wDB, tx: nil}, nil
}

type gormTxAdapter struct {
	wrapperDB *WrapperGorm
	tx        *gorm.DB // for ManualComplete
}

func (adapter *gormTxAdapter) AutoComplete(ctx context.Context, fn func(txCtx context.Context) error) error {
	if ctx == nil {
		ctx = context.Background()
	}

	if adapter.wrapperDB.ExistTxInsideContext(ctx) {
		return fn(ctx)
	}

	gormTxFn := func(tx *gorm.DB) error {
		txCtx := adapter.wrapperDB.NewTxContext(ctx, tx)
		return fn(txCtx)
	}

	err := adapter.wrapperDB.Unwrap().Transaction(gormTxFn)
	if err != nil {
		if errorY.IsUndefinedError(err) {
			return errorY.Wrap(errorY.ErrSystem, err.Error())
		}
		return err
	}

	return nil
}

func (adapter *gormTxAdapter) ManualComplete(
	ctx context.Context,
	fn func(txCtx context.Context) error,
) (
	commit func() error,
	rollback func() error,
	fnErr error,
) {
	if ctx == nil {
		ctx = context.Background()
	}

	if adapter.wrapperDB.ExistTxInsideContext(ctx) {
		return adapter.doNothing, adapter.doNothing, fn(ctx)
	}

	adapter.tx = adapter.wrapperDB.Unwrap().Begin()
	if err := adapter.tx.Error; err != nil {
		return adapter.doNothing, adapter.doNothing, errorY.Wrap(errorY.ErrSystem, err.Error())
	}

	txCtx := adapter.wrapperDB.NewTxContext(ctx, adapter.tx)
	return adapter.commit, adapter.rollback, fn(txCtx)
}

func (adapter *gormTxAdapter) doNothing() error { return nil }

func (adapter *gormTxAdapter) commit() error {
	err := adapter.tx.Commit().Error
	if err != nil {
		return errorY.Wrap(errorY.ErrSystem, err.Error())
	}
	return nil
}

func (adapter *gormTxAdapter) rollback() error {
	err := adapter.tx.Rollback().Error
	if err != nil {
		return errorY.Wrap(errorY.ErrSystem, err.Error())
	}
	return nil
}
