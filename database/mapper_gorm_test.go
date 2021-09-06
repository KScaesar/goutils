package database

import (
	"testing"

	"github.com/fatih/structs"
	"github.com/stretchr/testify/assert"
)

func TestTransformQueryParamToGorm(t *testing.T) {
	boolFalse := false
	IntZero := 0

	type Embed struct {
		Age int `rdb:"age = ?"`
	}

	tests := []struct {
		name        string
		param       interface{}
		expectedSql string
	}{
		{
			name: "struct is embed",
			param: struct {
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
			param: struct {
				OccurredAt string `rdb:"occurred_at is ?"`
			}{
				OccurredAt: "null",
			},
			expectedSql: "SELECT * FROM `book` WHERE occurred_at is ?",
		},
		{
			name: "string have a value",
			param: struct {
				Name string `rdb:"name = ?"`
			}{
				Name: "haha",
			},
			expectedSql: "SELECT * FROM `book` WHERE name = ?",
		},
		{
			name: "string is empty",
			param: struct {
				Name string `rdb:"name = ?"`
			}{
				Name: "",
			},
			expectedSql: "SELECT * FROM `book`",
		},
		{
			name: "bool is pointer and nil",
			param: struct {
				IsAdmin *bool `rdb:"is_admin = ?"`
			}{
				IsAdmin: nil,
			},
			expectedSql: "SELECT * FROM `book`",
		},
		{
			name: "bool is pointer and false",
			param: struct {
				IsAdmin *bool `rdb:"is_admin = ?"`
			}{
				IsAdmin: &boolFalse,
			},
			expectedSql: "SELECT * FROM `book` WHERE is_admin = ?",
		},
		{
			name: "bool is false",
			param: struct {
				IsAdmin bool `rdb:"is_admin = ?"`
			}{
				IsAdmin: false,
			},
			expectedSql: "SELECT * FROM `book`",
		},
		{
			name: "int is pointer and nil",
			param: struct {
				Age *int `rdb:"age = ?"`
			}{
				Age: nil,
			},
			expectedSql: "SELECT * FROM `book`",
		},
		{
			name: "search int == 0, need to int type is pointer",
			param: struct {
				Age *int `rdb:"age = ?"`
			}{
				Age: &IntZero,
			},
			expectedSql: "SELECT * FROM `book` WHERE age = ?",
		},
		{
			name: "int is zero",
			param: struct {
				Age int `rdb:"age = ?"`
			}{
				Age: 0,
			},
			expectedSql: "SELECT * FROM `book`",
		},
		{
			name: "id slice not nil",
			param: struct {
				IDSet []int `rdb:"id In (?)"`
			}{
				IDSet: []int{2, 4, 5},
			},
			expectedSql: "SELECT * FROM `book` WHERE id In (?,?,?)",
		},
		{
			name: "id slice is nil",
			param: struct {
				IDSet []int `rdb:"id In (?)"`
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

	db := mockGorm(false)

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			where := TransformQueryParamToGorm(tt.param)

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

	diff := UpdatedValue(man.Before, man)

	db := mockGorm(true)
	result := db.Table("person").Where("id = ?", man.ID).Updates(diff)
	assert.NoError(t, result.Error)

	actualSql := result.Statement.SQL.String()
	expected := "UPDATE `person` SET `IsAdmin`=?,`Money`=?,`Name`=? WHERE id = ?"
	assert.Equal(t, expected, actualSql)
}
