package database

import (
	"context"

	"gorm.io/gorm"

	"github.com/KScaesar/goutils/errors"
)

func NewGormTxFactory(db *WrapperGorm) TxFactory {
	return &gormTxFactory{wDB: db}
}

type gormTxFactory struct {
	wDB *WrapperGorm
}

func (f *gormTxFactory) CreateTx(ctx context.Context) Transaction {
	if ctx == nil {
		ctx = context.Background()
	}

	return &gormTxAdapter{wrapperDB: f.wDB, ctx: ctx}
}

type gormTxAdapter struct {
	wrapperDB *WrapperGorm

	tx  *gorm.DB
	ctx context.Context
}

func (adapter *gormTxAdapter) AutoComplete(fn func(txCtx context.Context) error) error {
	ctx := adapter.ctx

	if adapter.wrapperDB.ExistTxInsideContext(ctx) {
		return fn(ctx)
	}

	gormTxFn := func(tx *gorm.DB) error {
		txCtx := adapter.wrapperDB.ContextWithTx(ctx, tx)
		return fn(txCtx)
	}

	err := adapter.wrapperDB.Unwrap().Transaction(gormTxFn)
	if err != nil {
		if errors.IsUndefinedError(err) {
			return errors.Wrap(errors.ErrSystem, err.Error())
		}
		return err
	}

	return nil
}

func (adapter *gormTxAdapter) ManualComplete(fn func(txCtx context.Context) error) (
	commit func() error,
	rollback func() error,
	fnErr error,
) {
	ctx := adapter.ctx

	if adapter.wrapperDB.ExistTxInsideContext(ctx) {
		return adapter.doNothing, adapter.doNothing, fn(ctx)
	}

	adapter.tx = adapter.wrapperDB.Unwrap().Begin()
	if err := adapter.tx.Error; err != nil {
		return adapter.doNothing, adapter.doNothing, errors.Wrap(errors.ErrSystem, err.Error())
	}

	txCtx := adapter.wrapperDB.ContextWithTx(ctx, adapter.tx)
	return adapter.commit, adapter.rollback, fn(txCtx)
}

func (adapter *gormTxAdapter) doNothing() error { return nil }

func (adapter *gormTxAdapter) commit() error {
	err := adapter.tx.Commit().Error
	if err != nil {
		return errors.Wrap(errors.ErrSystem, err.Error())
	}
	return nil
}

func (adapter *gormTxAdapter) rollback() error {
	err := adapter.tx.Rollback().Error
	if err != nil {
		return errors.Wrap(errors.ErrSystem, err.Error())
	}
	return nil
}
