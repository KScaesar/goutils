//go:build integration
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
	db := pgGorm(nil)

	tx := db.Unwrap().Begin()
	defer tx.Commit()
	txCtx := db.ContextWithTx(nil, tx)
	actualProcessor := db.ChooseProcessor(txCtx)

	assert.Equal(t, tx, actualProcessor)
	assert.NotEqual(t, db.Unwrap(), actualProcessor)
}

//go:generate go test -v -tags=integration -run TestWrapperGorm_EnableTransaction_but_different_database
func TestWrapperGorm_EnableTransaction_but_different_database(t *testing.T) {
	db1 := pgGorm(nil)
	db2 := pgGorm(nil)

	tx1 := db1.Unwrap().Begin()
	defer tx1.Commit()
	txCtx1 := db1.ContextWithTx(nil, tx1)
	actualProcessor := db2.ChooseProcessor(txCtx1)

	assert.Equal(t, db2.Unwrap(), actualProcessor)
	assert.NotEqual(t, tx1, actualProcessor)
}

func TestWrapperGorm_NoTransaction_when_ctx_is_nil(t *testing.T) {
	db := pgGorm(nil)

	actualProcessor := db.ChooseProcessor(nil)

	assert.Equal(t, db.Unwrap(), actualProcessor)
}

func TestWrapperGorm_NoTransaction_when_ctx_not_nil(t *testing.T) {
	db := pgGorm(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	actualProcessor := db.ChooseProcessor(ctx)

	assert.Equal(t, db.Unwrap(), actualProcessor)
}
