package database_test

import (
	"testing"

	"github.com/fatih/structs"
	"github.com/stretchr/testify/assert"

	"github.com/KScaesar/goutils"
	"github.com/KScaesar/goutils/database"
	"github.com/KScaesar/goutils/xTest"
)

func TestGormFilter(t *testing.T) {
	boolFalse := false
	IntZero := 0

	type Embed struct {
		Age int `rdb:"age = ?"`
	}

	tests := []struct {
		name        string
		filter      interface{}
		expectedSql string
	}{
		{
			name: "struct is embed",
			filter: struct {
				Embed
				IsAdmin *bool `rdb:"is_admin = ?"`
			}{
				IsAdmin: &boolFalse,
				Embed: Embed{
					Age: 30,
				},
			},
			expectedSql: "SELECT * FROM `book` WHERE age = ? AND is_admin = ?",
		},
		{
			name: "search null",
			filter: struct {
				OccurredAt string `rdb:"occurred_at is ?"`
			}{
				OccurredAt: "null",
			},
			expectedSql: "SELECT * FROM `book` WHERE occurred_at is ?",
		},
		{
			name: "string have a value",
			filter: struct {
				Name string `rdb:"name = ?"`
			}{
				Name: "haha",
			},
			expectedSql: "SELECT * FROM `book` WHERE name = ?",
		},
		{
			name: "string is empty",
			filter: struct {
				Name string `rdb:"name = ?"`
			}{
				Name: "",
			},
			expectedSql: "SELECT * FROM `book`",
		},
		{
			name: "bool is pointer and nil",
			filter: struct {
				IsAdmin *bool `rdb:"is_admin = ?"`
			}{
				IsAdmin: nil,
			},
			expectedSql: "SELECT * FROM `book`",
		},
		{
			name: "bool is pointer and false",
			filter: struct {
				IsAdmin *bool `rdb:"is_admin = ?"`
			}{
				IsAdmin: &boolFalse,
			},
			expectedSql: "SELECT * FROM `book` WHERE is_admin = ?",
		},
		{
			name: "bool is false",
			filter: struct {
				IsAdmin bool `rdb:"is_admin = ?"`
			}{
				IsAdmin: false,
			},
			expectedSql: "SELECT * FROM `book`",
		},
		{
			name: "int is pointer and nil",
			filter: struct {
				Age *int `rdb:"age = ?"`
			}{
				Age: nil,
			},
			expectedSql: "SELECT * FROM `book`",
		},
		{
			name: "search int == 0, need to int type is pointer",
			filter: struct {
				Age *int `rdb:"age = ?"`
			}{
				Age: &IntZero,
			},
			expectedSql: "SELECT * FROM `book` WHERE age = ?",
		},
		{
			name: "int is zero",
			filter: struct {
				Age int `rdb:"age = ?"`
			}{
				Age: 0,
			},
			expectedSql: "SELECT * FROM `book`",
		},
		{
			name: "id slice not nil",
			filter: struct {
				IDSet []int `rdb:"id In (?)"`
			}{
				IDSet: []int{2, 4, 5},
			},
			expectedSql: "SELECT * FROM `book` WHERE id In (?,?,?)",
		},
		{
			name: "id slice is nil",
			filter: struct {
				IDSet []int `rdb:"id In (?)"`
			}{
				IDSet: nil,
			},
			expectedSql: "SELECT * FROM `book`",
		},
		{
			name:        "filter struct is nil",
			filter:      nil,
			expectedSql: "SELECT * FROM `book`",
		},
		{
			name: "FilterOption type",
			filter: goutils.FilterOption{
				{"user_id = ?", 123},
				{"datetime in (?)", []string{"2020-10-17", "2021-09-16"}},
			},
			expectedSql: "SELECT * FROM `book` WHERE user_id = ? AND datetime in (?,?)",
		},
	}

	db := xTest.MockGormMysql(false).Unwrap()

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			where := database.GormFilter(tt.filter)

			var data []map[string]interface{}
			stmt := db.Table("book").Scopes(where...).Find(&data).Statement
			actualSql := stmt.SQL.String()

			assert.Equal(t, tt.expectedSql, actualSql)
		})
	}
}

func TestUpdatedValue(t *testing.T) {
	type Person struct {
		ID      int
		Name    string
		Money   int
		Email   string
		IsAdmin bool

		Before map[string]interface{} `gorm:"-"`
	}

	man := &Person{
		ID:      3,
		Name:    "caesar",
		Money:   38,
		Email:   "x246libra@hotmail.com",
		IsAdmin: true,
	}
	man.Before = structs.New(man).Map()

	man.Name = "peter"
	man.Money = 0
	man.IsAdmin = false

	diff := database.UpdatedValue(man.Before, man)

	db := xTest.MockGormMysql(true).Unwrap()
	result := db.Table("person").Where("id = ?", man.ID).Updates(diff)
	assert.NoError(t, result.Error)

	actualSql := result.Statement.SQL.String()
	expected := "UPDATE `person` SET `IsAdmin`=?,`Money`=?,`Name`=? WHERE id = ?"
	assert.Equal(t, expected, actualSql)
}
