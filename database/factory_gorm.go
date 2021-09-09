package database

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormMysql(cfg *RMDBConfig, debug bool) (*WrapperGorm, error) {
	cfg.setDefaultValue()

	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?parseTime=true&loc=Local&charset=utf8mb4",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxConn)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)

	gormDB, err := gorm.Open(
		mysql.New(mysql.Config{Conn: sqlDB}),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}

	if debug {
		gormDB = gormDB.Debug()
	}

	return &WrapperGorm{gormDB}, nil
}

func NewGormPostgres(cfg *RMDBConfig, debug bool) (*WrapperGorm, error) {
	cfg.setDefaultValue()

	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s lock_timeout=5000 idle_in_transaction_session_timeout=10000 sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	gormDB, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxConn)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)

	if debug {
		gormDB = gormDB.Debug()
	}

	return &WrapperGorm{gormDB}, nil
}
