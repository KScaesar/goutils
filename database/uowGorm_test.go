// +build integration

package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Min-Feng/goutils/database"
)

func Test_uowGorm_AutoStart(t *testing.T) {
	fixture := testFixture{}
	db := fixture.pgGorm()

	// https://gorm.io/docs/migration.html#Tables
	sqlBook := &infraBook{}
	if !db.Unwrap().Migrator().HasTable(sqlBook.TableName()) {
		assert.NoError(t, db.Unwrap().Migrator().CreateTable(sqlBook))
	}

	uowFactory := database.NewUowGormFactory(db)
	repo := bookGormRepo{db: db, tableName: sqlBook.TableName()}

	uow, err := uowFactory.CreateUow()
	assert.NoError(t, err)

	createFn := func(name string) func(txCtx context.Context) error {
		return func(txCtx context.Context) error {
			book := &DomainBook{
				SqlID:    uuid.New().String(),
				Name:     "ddd_is_good" + "#" + name,
				NoTzTime: time.Now(),
				TzTime:   time.Now(),
			}
			if err := repo.createBook(txCtx, book); err != nil {
				return err
			}

			book.Name = "tdd_is_good" + "#" + name
			if err := repo.updateBook(txCtx, book); err != nil {
				return err
			}

			return nil
		}
	}

	txFn := createFn("tx")
	uowErr := uow.AutoStart(nil, txFn)
	assert.NoError(t, uowErr, "enable tx")

	noTxFn := createFn("noTx")
	fnErr := noTxFn(nil)
	assert.NoError(t, fnErr, "not enable transaction")
}

type bookGormRepo struct {
	db        *database.WrapperGorm
	tableName string
}

func (repo *bookGormRepo) createBook(ctx context.Context, book *DomainBook) error {
	p := repo.db.GetTxFromCtxAndSelectProcessor(ctx)
	return p.Table(repo.tableName).Create(book).Error
}

func (repo *bookGormRepo) updateBook(ctx context.Context, book *DomainBook) error {
	p := repo.db.GetTxFromCtxAndSelectProcessor(ctx)
	return p.Table(repo.tableName).Model(book).Updates(book).Error
}
