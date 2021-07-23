package database

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewGormMysql(cfg *RMDBConfig) (*WrapperGorm, error) {
	sqlDB, err := sql.Open("mysql", cfg.MysqlDSN())
	if err != nil {
		return nil, err
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDB.SetMaxOpenConns(cfg.MaxConn)

	gormDB, err := gorm.Open(
		mysql.New(mysql.Config{Conn: sqlDB}),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}

	return &WrapperGorm{gormDB}, nil
}

type WrapperGorm struct {
	db *gorm.DB
}

// NewTxContext 如果 ctx 有數值, 此 method 會自動和 tx gorm.DB 進行關聯
func (wrapper *WrapperGorm) NewTxContext(ctx context.Context, tx *gorm.DB) (txCtx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	// key 用來確認是否來自同一個 *gorm.DB
	key := wrapper
	return context.WithValue(ctx, key, tx.WithContext(ctx))
}

// GetTxFromCtxAndSelectProcessor
// 如果找不到符合的 tx 元件, 則使用原本的 db 元件,
// 找不到的原因:
//
// 1. 沒有 transaction 需求, 所以沒有傳入 tx 元件
// 2. tx 元件來自不同的 database, 那當然無法達成 transaction 的需求, 只能各自處理 sql operation
//
// 如果 txCtx 有數值, 此 method 會自動和 db gorm.DB 進行關聯
//
// 回傳值 processor, 可以代表 tx or db, 兩者型別都是  *gorm.DB
// 我覺得 gorm 這設計不太好, 不同用途放在同一個具體型別
// 其差異請參考下列網址
// https://gorm.io/docs/method_chaining.html#New-Session-Mode
func (wrapper *WrapperGorm) GetTxFromCtxAndSelectProcessor(txCtx context.Context) (processor *gorm.DB) {
	if txCtx == nil {
		return wrapper.db
	}

	// key 用來確認是否來自同一個 *gorm.DB
	key := wrapper
	tx, ok := txCtx.Value(key).(*gorm.DB)
	if ok {
		return tx // NewTxContext 的時候, 已經執行過 gorm.DB WithContext
	}
	return wrapper.db.WithContext(txCtx)
}

func (wrapper *WrapperGorm) Unwrap() *gorm.DB {
	return wrapper.db
}
