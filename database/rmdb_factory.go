package database

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
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

func NewGormPostgres(cfg *RMDBConfig) (*WrapperGorm, error) {
	gormDB, err := gorm.Open(
		postgres.Open(cfg.GormPgDSN()),
		&gorm.Config{},
	)
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxConn_())
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn_())

	return &WrapperGorm{gormDB}, nil
}

type RMDBConfig struct {
	User        string
	Password    string
	Host        string
	Port        string
	Database    string
	MaxConn     int
	MaxIdleConn int
}

func (c *RMDBConfig) MysqlDSN() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true&loc=Local&charset=utf8mb4", c.User, c.Password, c.Host, c.Port, c.Database)
}

func (c *RMDBConfig) GormPgDSN() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Database)
}

func (c *RMDBConfig) MaxConn_() int {
	if c.MaxConn <= 0 {
		const defaultSize = 8
		return defaultSize
	}
	return c.MaxConn
}

func (c *RMDBConfig) MaxIdleConn_() int {
	if c.MaxConn <= 0 {
		const defaultSize = 4
		return defaultSize
	}
	return c.MaxIdleConn
}
