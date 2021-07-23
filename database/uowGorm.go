package database

import (
	"context"

	"gorm.io/gorm"

	"github.com/Min-Feng/goutils/errorY"
)

func NewUowGormFactory(db *WrapperGorm) UowFactory {
	return &uowGormFactory{wDB: db}
}

type uowGormFactory struct {
	wDB *WrapperGorm
}

func (f *uowGormFactory) CreateUow() (Uow, error) {
	return &uowGorm{wrapperDB: f.wDB, tx: nil}, nil
}

type uowGorm struct {
	wrapperDB *WrapperGorm
	tx        *gorm.DB // for manual start
}

func (uow *uowGorm) AutoStart(ctx context.Context, txFn func(txCtx context.Context) error) error {
	if ctx == nil {
		ctx = context.Background()
	}

	gormTxFn := func(tx *gorm.DB) error {
		txCtx := uow.wrapperDB.NewTxContext(ctx, tx)
		return txFn(txCtx)
	}

	err := uow.wrapperDB.Unwrap().Transaction(gormTxFn)
	if err != nil {
		if errorY.IsUndefinedError(err) {
			return errorY.Wrap(errorY.ErrSystem, err.Error())
		}
		return err
	}

	return nil
}

func (uow *uowGorm) ManualStart(ctx context.Context, txFn func(txCtx context.Context) error) error {
	if ctx == nil {
		ctx = context.Background()
	}

	uow.tx = uow.wrapperDB.Unwrap().Begin()
	if err := uow.tx.Error; err != nil {
		return errorY.Wrap(errorY.ErrSystem, err.Error())
	}

	txCtx := uow.wrapperDB.NewTxContext(ctx, uow.tx)
	return txFn(txCtx)
}

func (uow *uowGorm) Commit() error {
	if uow.tx == nil {
		return errorY.Wrap(errorY.ErrSystem, "invalid command, need to call ManualStart")
	}

	err := uow.tx.Commit().Error
	if err != nil {
		return errorY.Wrap(errorY.ErrSystem, err.Error())
	}
	return nil
}

func (uow *uowGorm) Rollback() error {
	if uow.tx == nil {
		return errorY.Wrap(errorY.ErrSystem, "invalid command, need to call ManualStart")
	}

	err := uow.tx.Rollback().Error
	if err != nil {
		return errorY.Wrap(errorY.ErrSystem, err.Error())
	}
	return nil
}
