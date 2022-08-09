//go:build integration
// +build integration

package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/KScaesar/goutils"
	"github.com/KScaesar/goutils/database"
	"github.com/KScaesar/goutils/xLog"
)

func Test_txGorm_AutoComplete(t *testing.T) {
	db := pgGorm(nil)

	// https://gorm.io/docs/migration.html#Tables
	tableName := "testing_books"
	if !db.Unwrap().Migrator().HasTable(tableName) {
		err := db.Unwrap().
			Table(tableName).
			Migrator().
			CreateTable(&DomainBook{})
		assert.NoError(t, err)
	}

	txFactory := database.NewGormTxFactory(db)
	repo := bookGormRepo{db: db, tableName: tableName}

	tx := txFactory.CreateTx(nil)

	id := goutils.NewULID()
	fn := func(name string) func(txCtx context.Context) error {
		return func(txCtx context.Context) error {
			now := time.Now()
			book := &DomainBook{
				SqlID:    id,
				Name:     "python" + "#" + name,
				NoTzTime: now,
				TzTime:   now,
				NullTime: goutils.Time(now),
				AutoTIme: goutils.Time(now),
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

	err := tx.AutoComplete(fn("enable tx"))
	assert.NoError(t, err, "enable tx")

	// err = fn("disable tx")(nil)
	// assert.NoError(t, err, "disable tx")

	book, err := repo.getBook(nil, id)
	assert.NoError(t, err)
	xLog.Info().Interface("book", book).Send()
}

type bookGormRepo struct {
	db        *database.WrapperGorm
	tableName string
}

func (repo *bookGormRepo) createBook(ctx context.Context, book *DomainBook) error {
	p := repo.db.SelectProcessor(ctx)
	return p.Table(repo.tableName).Create(book).Error
}

func (repo *bookGormRepo) updateBook(ctx context.Context, book *DomainBook) error {
	p := repo.db.SelectProcessor(ctx)
	return p.Table(repo.tableName).Save(book).Error
}

func (repo *bookGormRepo) getBook(ctx context.Context, id string) (book DomainBook, err error) {
	p := repo.db.SelectProcessor(ctx)
	return book, p.Table(repo.tableName).Where("id = ?", id).Find(&book).Error
}
