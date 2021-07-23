// +build integration

package database_test

import "github.com/Min-Feng/goutils/database"

type testFixture struct{}

func (f testFixture) mysqlConnectConfig() database.RMDBConfig {
	return database.RMDBConfig{
		User:        "root",
		Password:    "1234",
		Host:        "localhost",
		Port:        "3306",
		Database:    "GoCodeStyle",
		MaxConn:     8,
		MaxIdleConn: 4,
	}
}

func (f testFixture) gormMysql(cfg database.RMDBConfig) *database.WrapperGorm {
	db, err := database.NewGormMysql(&cfg)
	if err != nil {
		panic(err)
	}
	return db
}
