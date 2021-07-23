// +build integration

package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//go:generate go test -v -tags=integration -run TestWrapperGorm_EnableTransaction
func TestWrapperGorm_EnableTransaction(t *testing.T) {
	fixture := testFixture{}
	db := fixture.gormMysql(fixture.mysqlConnectConfig())

	tx := db.Unwrap().Begin()
	defer tx.Commit()
	txCtx := db.NewTxContext(nil, tx)
	actualProcessor := db.GetTxFromCtxAndSelectProcessor(txCtx)

	expectedProcessor := tx
	assert.Equal(t, expectedProcessor, actualProcessor)

	notExpectedProcessor := db.Unwrap()
	assert.NotEqual(t, notExpectedProcessor, actualProcessor)
}

//go:generate go test -v -tags=integration -run TestWrapperGorm_EnableTransaction_but_different_database
func TestWrapperGorm_EnableTransaction_but_different_database(t *testing.T) {
	fixture := testFixture{}
	cfg := fixture.mysqlConnectConfig()
	db1 := fixture.gormMysql(cfg)
	db2 := fixture.gormMysql(cfg)

	tx1 := db1.Unwrap().Begin()
	defer tx1.Commit()
	txCtx1 := db1.NewTxContext(nil, tx1)
	actualProcessor := db2.GetTxFromCtxAndSelectProcessor(txCtx1)

	expectedProcessor := db2.Unwrap()
	assert.Equal(t, expectedProcessor, actualProcessor)

	notExpectedProcessor := tx1
	assert.NotEqual(t, notExpectedProcessor, actualProcessor)
}

func TestWrapperGorm_NoTransaction_when_ctx_is_nil(t *testing.T) {
	fixture := testFixture{}
	db := fixture.gormMysql(fixture.mysqlConnectConfig())

	actualProcessor := db.GetTxFromCtxAndSelectProcessor(nil)

	expectedProcessor := db.Unwrap()
	assert.Equal(t, expectedProcessor, actualProcessor)
}

func TestWrapperGorm_NoTransaction_when_ctx_not_nil(t *testing.T) {
	fixture := testFixture{}
	db := fixture.gormMysql(fixture.mysqlConnectConfig())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	actualProcessor := db.GetTxFromCtxAndSelectProcessor(ctx)

	expectedProcessor := db.Unwrap()
	assert.Equal(t, expectedProcessor, actualProcessor)
}
