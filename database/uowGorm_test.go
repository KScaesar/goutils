// +build integration

package database_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Min-Feng/goutils/database"
)

func Test_uowGorm_AutoStart(t *testing.T) {
	fixture := testFixture{}
	db := fixture.mysqlGorm(fixture.mysqlConnectConfig())

	// https://gorm.io/docs/migration.html#Tables
	sqlBookConfig := &sqlBook{}
	if !db.Unwrap().Migrator().HasTable(sqlBookConfig.TableName()) {
		db.Unwrap().Migrator().CreateTable(sqlBookConfig)
	}

	uowFactory := database.NewUowGormFactory(db)
	repo := bookRepo{db: db, tableName: sqlBookConfig.TableName()}

	uow, err := uowFactory.CreateUow()
	assert.NoError(t, err)

	fn := func(txCtx context.Context) error {
		book1 := &DomainBook{Name: "ddd_is_good"}
		if err := repo.createBook(txCtx, book1); err != nil {
			return err
		}

		book1.Name = "tdd_is_good"
		if err := repo.updateBook(txCtx, book1); err != nil {
			return err
		}

		return nil
	}

	// enable tx
	uowErr := uow.AutoStart(nil, fn)
	assert.NoError(t, uowErr)

	// normal sql, not enable transaction
	// _ = uow
	// repoErr := fn(nil)
	// assert.NoError(t, repoErr)
}

type sqlBook struct {
	DomainBook
}

func (s *sqlBook) TableName() string {
	return "testing_books"
}

type bookRepo struct {
	db        *database.WrapperGorm
	tableName string
}

func (repo *bookRepo) createBook(ctx context.Context, book *DomainBook) error {
	p := repo.db.GetTxFromCtxAndSelectProcessor(ctx)
	return p.Table(repo.tableName).Create(book).Error
}

func (repo *bookRepo) updateBook(ctx context.Context, book *DomainBook) error {
	p := repo.db.GetTxFromCtxAndSelectProcessor(ctx)
	return p.Table(repo.tableName).Model(book).Updates(book).Error
}
