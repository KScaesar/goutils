package database

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// mockGorm
//
// https://stackoverflow.com/questions/58804606/go-unit-tests-call-to-database-transaction-begin-was-not-expected-error
// https://gorm.io/docs/sql_builder.html#DryRun-Mode
func mockGorm(debug bool) *gorm.DB {
	sqlDB, _, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(
		mysql.New(mysql.Config{
			Conn:                      sqlDB,
			SkipInitializeWithVersion: true,
		}),
		&gorm.Config{
			SkipDefaultTransaction: true,
			DryRun:                 true,
		},
	)
	if err != nil {
		panic(err)
	}

	if debug {
		return db.Debug()
	}

	return db
}
