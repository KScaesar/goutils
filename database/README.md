# introduction

實現 uow(Unit Of Work) 功能  
或者說 資料庫交易功能

使用方式, 可以查看 test function 得知

## interface define

[link](./transaction.go)

```go
package database

type TxFactory interface {
	CreateTx(ctx context.Context) Transaction
}

type Transaction interface {
	// AutoComplete 執行 fn 之後, if fn return nil 自動執行 Commit, else return err 自動執行 Rollback
	AutoComplete(fn func(txCtx context.Context) error) error

	// ManualComplete 執行 fn 之後, 需要另外手動呼叫 Commit or Rollback
	ManualComplete(fn func(txCtx context.Context) error) (
		commit func() error,
		rollback func() error,
		fnErr error,
	)
}
```

## implement

1. [gorm](./txGorm.go)
2. [mongo DB](./txMongo.go)
