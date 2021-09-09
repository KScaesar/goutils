package database

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func mysqlDSN(c *RMDBConfig) string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&loc=Local&charset=utf8mb4", c.User, c.Password, c.Host, c.Port, c.Database)
}

func NewGormMysql(cfg *RMDBConfig) (*WrapperGorm, error) {
	sqlDB, err := sql.Open("mysql", mysqlDSN(cfg))
	if err != nil {
		return nil, err
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxConn_())
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn_())

	gormDB, err := gorm.Open(
		mysql.New(mysql.Config{Conn: sqlDB}),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}

	return &WrapperGorm{gormDB}, nil
}

func gormPgDSN(c *RMDBConfig) string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s lock_timeout=5000 idle_in_transaction_session_timeout=10000 sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Database)
}

func NewGormPostgres(cfg *RMDBConfig, debug bool) (*WrapperGorm, error) {
	gormDB, err := gorm.Open(
		postgres.Open(gormPgDSN(cfg)),
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

	sqlDB.SetMaxOpenConns(cfg.MaxConn_())
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn_())

	if debug {
		gormDB = gormDB.Debug()
	}

	return &WrapperGorm{gormDB}, nil
}
