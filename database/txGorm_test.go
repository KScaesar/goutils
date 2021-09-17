// +build integration

package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Min-Feng/goutils/database"
)

func Test_txGorm_AutoComplete(t *testing.T) {
	db := pgGorm(nil)

	// https://gorm.io/docs/migration.html#Tables
	sqlBook := &infraBook{}
	if !db.Unwrap().Migrator().HasTable(sqlBook.TableName()) {
		assert.NoError(t, db.Unwrap().Migrator().CreateTable(sqlBook))
	}

	txFactory := database.NewGormTxFactory(db)
	repo := bookGormRepo{db: db, tableName: sqlBook.TableName()}

	tx, err := txFactory.CreateTx(nil)
	assert.NoError(t, err)

	fn := func(name string) func(txCtx context.Context) error {
		return func(txCtx context.Context) error {
			book := &DomainBook{
				Name:     "python" + "#" + name,
				NoTzTime: time.Now(),
				TzTime:   time.Now(),
			}
			if err := repo.createBook(txCtx, book); err != nil {
				return err
			}

			book.Name = "golang" + "#" + name
			if err := repo.updateBook(txCtx, book); err != nil {
				return err
			}

			return nil
		}
	}

	err = tx.AutoComplete(fn("tx"))
	assert.NoError(t, err, "enable tx")

	err = fn("noTx")(nil)
	assert.NoError(t, err, "not enable tx")
}

type bookGormRepo struct {
	db        *database.WrapperGorm
	tableName string
}

func (repo *bookGormRepo) createBook(ctx context.Context, book *DomainBook) error {
	p := repo.db.TxFromContextAndSelectProcessor(ctx)
	return p.Table(repo.tableName).Create(book).Error
}

func (repo *bookGormRepo) updateBook(ctx context.Context, book *DomainBook) error {
	p := repo.db.TxFromContextAndSelectProcessor(ctx)
	return p.Table(repo.tableName).Save(book).Error
}
