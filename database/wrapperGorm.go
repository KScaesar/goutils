package database

import (
	"context"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type WrapperGorm struct {
	db *gorm.DB
}

func (wrapper *WrapperGorm) NewTxContext(ctx context.Context, tx *gorm.DB) (txCtx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	// key 用來確認是否來自同一個 *gorm.DB
	key := wrapper
	return context.WithValue(ctx, key, tx)
}

// GetTxFromCtxAndSelectProcessor
// 如果找不到符合的 tx 元件, 則使用原本的 db 元件,
// 找不到的原因:
//
// 1. 沒有 transaction 需求, 所以沒有傳入 tx 元件
// 2. tx 元件來自不同的 database, 那當然無法達成 transaction 的需求, 只能各自處理 sql operation
//
// 回傳值 processor, 可以代表 tx or db, 兩者型別都是  *gorm.DB
// 我覺得 gorm 這設計不太好, 不同用途放在同一個具體型別
// 其差異請參考下列網址
// https://gorm.io/docs/method_chaining.html#New-Session-Mode
func (wrapper *WrapperGorm) GetTxFromCtxAndSelectProcessor(txCtx context.Context) (processor *gorm.DB) {
	if txCtx == nil {
		return wrapper.Unwrap()
	}

	// key 用來確認是否來自同一個 *gorm.DB
	key := wrapper
	tx, ok := txCtx.Value(key).(*gorm.DB)
	if ok {
		return tx // NewTxContext 的時候, 已經執行過 gorm.DB WithContext
	}
	return wrapper.Unwrap()
}

func (wrapper *WrapperGorm) Unwrap() *gorm.DB {
	return wrapper.db
}
