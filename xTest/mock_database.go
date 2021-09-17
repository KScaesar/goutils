package xTest

import (
	"context"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Min-Feng/goutils/database"
)

// MockGormMysql
//
// https://stackoverflow.com/questions/58804606/go-unit-tests-call-to-database-transaction-begin-was-not-expected-error
// https://gorm.io/docs/sql_builder.html#DryRun-Mode
func MockGormMysql(debug bool) *database.WrapperGorm {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	mock.ExpectBegin()
	mock.ExpectCommit()

	db, err := gorm.Open(
		mysql.New(
			mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			},
		),
		&gorm.Config{
			DryRun: true,
		},
	)
	if err != nil {
		panic(err)
	}

	if debug {
		db = db.Debug()
	}

	return database.NewWrapperGorm(db)
}

func MockGormPgsql(debug bool) *database.WrapperGorm {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	mock.ExpectBegin()
	mock.ExpectCommit()

	db, err := gorm.Open(
		postgres.New(
			postgres.Config{
				Conn: sqlDB,
			},
		),
		&gorm.Config{
			DryRun: true,
		},
	)
	if err != nil {
		panic(err)
	}

	if debug {
		db = db.Debug()
	}

	return database.NewWrapperGorm(db)
}

type MockTxFactory struct{}

func (MockTxFactory) CreateTx(ctx context.Context) (database.Transaction, error) {
	spy := &spyTxAdapter{ctx: ctx}
	return spy, nil
}

type spyTxAdapter struct {
	ctx                   context.Context
	executeAutoComplete   bool
	executeManualComplete bool
	executeCommit         bool
	executeRollback       bool
}

func (adapter *spyTxAdapter) AutoComplete(fn func(txCtx context.Context) error) error {
	adapter.executeAutoComplete = true
	return fn(adapter.ctx)
}

func (adapter *spyTxAdapter) ManualComplete(fn func(txCtx context.Context) error) (
	commit func() error,
	rollback func() error,
	fnErr error,
) {
	adapter.executeManualComplete = true
	commit = func() error {
		adapter.executeCommit = true
		return nil
	}
	rollback = func() error {
		adapter.executeRollback = true
		return nil
	}
	return commit, rollback, fn(adapter.ctx)
}

func DoesTxAutoComplete(tx database.Transaction) bool {
	return tx.(*spyTxAdapter).executeAutoComplete
}

func DoesTxManualComplete(tx database.Transaction) bool {
	return tx.(*spyTxAdapter).executeManualComplete
}

func DoesTxCommit(tx database.Transaction) bool {
	return tx.(*spyTxAdapter).executeCommit
}

func DoesTxRollback(tx database.Transaction) bool {
	return tx.(*spyTxAdapter).executeRollback
}
