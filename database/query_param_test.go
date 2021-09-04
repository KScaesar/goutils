package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestTransformQueryParamToGorm(t *testing.T) {
	boolFalse := false
	IntZero := 0

	type Embed struct {
		Age int `xQuery:"age = ?"`
	}

	tests := []struct {
		name        string
		param       interface{}
		expectedSql string
	}{
		{
			name: "search null",
			param: struct {
				OccurredAt string `xQuery:"occurred_at is ?"`
			}{
				OccurredAt: "null",
			},
			expectedSql: "SELECT * FROM `book` WHERE occurred_at is ?",
		},
		{
			name: "string have a value",
			param: struct {
				Name string `xQuery:"name = ?"`
			}{
				Name: "haha",
			},
			expectedSql: "SELECT * FROM `book` WHERE name = ?",
		},
		{
			name: "string is empty",
			param: struct {
				Name string `xQuery:"name = ?"`
			}{
				Name: "",
			},
			expectedSql: "SELECT * FROM `book`",
		},
		{
			name: "struct is embed",
			param: struct {
				IsAdmin bool `xQuery:"is_admin = ?"`
				Embed
			}{
				IsAdmin: false,
				Embed: Embed{
					Age: 30,
				},
			},
			expectedSql: "SELECT * FROM `book` WHERE age = ?",
		},
		{
			name: "bool is false",
			param: struct {
				IsAdmin bool `xQuery:"is_admin = ?"`
			}{
				IsAdmin: false,
			},
			expectedSql: "SELECT * FROM `book`",
		},
		{
			name: "bool is pointer and nil",
			param: struct {
				IsAdmin *bool `xQuery:"is_admin = ?"`
			}{
				IsAdmin: nil,
			},
			expectedSql: "SELECT * FROM `book`",
		},
		{
			name: "bool is pointer and false",
			param: struct {
				IsAdmin *bool `xQuery:"is_admin = ?"`
			}{
				IsAdmin: &boolFalse,
			},
			expectedSql: "SELECT * FROM `book` WHERE is_admin = ?",
		},
		{
			name: "int is pointer and nil",
			param: struct {
				Age *int `xQuery:"age = ?"`
			}{
				Age: nil,
			},
			expectedSql: "SELECT * FROM `book`",
		},
		{
			name: "search int == 0, need to int type is pointer",
			param: struct {
				Age *int `xQuery:"age = ?"`
			}{
				Age: &IntZero,
			},
			expectedSql: "SELECT * FROM `book` WHERE age = ?",
		},
		{
			name: "int is zero",
			param: struct {
				Age int `xQuery:"age = ?"`
			}{
				Age: 0,
			},
			expectedSql: "SELECT * FROM `book`",
		},
		{
			name: "id slice not nil",
			param: struct {
				IDSet []int `xQuery:"id In (?)"`
			}{
				IDSet: []int{2, 4, 5},
			},
			expectedSql: "SELECT * FROM `book` WHERE id In (?,?,?)",
		},
		{
			name: "id slice is nil",
			param: struct {
				IDSet []int `xQuery:"id In (?)"`
			}{
				IDSet: nil,
			},
			expectedSql: "SELECT * FROM `book`",
		},
		{
			name:        "param struct is nil",
			param:       nil,
			expectedSql: "SELECT * FROM `book`",
		},
	}

	sqlDB, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer sqlDB.Close()

	db, err := gorm.Open(
		mysql.New(mysql.Config{SkipInitializeWithVersion: true, Conn: sqlDB}),
		&gorm.Config{},
	)
	assert.NoError(t, err)

	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	db = db.Session(&gorm.Session{DryRun: true}).Debug()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			where := TransformQueryParamToGorm(tt.param)

			data := map[string]interface{}{}
			stmt := db.Table("book").Scopes(where...).Find(&data).Statement
			actualSql := stmt.SQL.String()

			assert.Equal(t, tt.expectedSql, actualSql)
		})
	}
}
