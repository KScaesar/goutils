package database

import (
	"reflect"

	"github.com/fatih/structs"
	"gorm.io/gorm"
)

// WhereParam 簡易查詢用途, 不可在 Value 嵌入其他條件
//
// example:
//
// query := WhereParam{
//     {"user_id = ?", 123},
//     {"datetime in (?)", []string{"2020-10-17", "2021-09-16"}},
//     {"is_admin = ?", true},
//     {"blocked_at = ?", "null"},
// }
type WhereParam []struct {
	Sql   string
	Value interface{}
}

// TransformWhereParamToGorm param 可傳入 struct or WhereParam,
// struct 會解析 tag key = rdb 的內容,
// 其 field 數值 若為 go 的初始值 會忽略此條件,
// 需要在 struct tag 填入 適當的 sql where 敘述.
//
// 若想查 xxx=false, 將型別定義為 pointer *bool;
// 若想查 xxx is null, 將型別定義為 string;
// 理論上全部數值可以都用 string 來表示.
//
// example:
//
// type CustomQuery struct {
//     UserIDs     []int     `rdb:"user_id in (?)"`
//     IsAdmin     *bool     `rdb:"is_admin = ?"`
//     BlockedTime time.Time `rdb:"blocked_time is ?"`
//     IsBlocked   string    `rdb:"is_blocked = ?"`
// }
//
// query := CustomQuery{
//     UserIDs:     []int{2, 4},
//     IsAdmin:     &BoolFalse,
//     BlockedTime: time.Now(),
//     IsBlocked:   "false",
// }
func TransformWhereParamToGorm(param interface{}) []func(db *gorm.DB) *gorm.DB {
	if param == nil {
		return nil
	}

	if reflect.TypeOf(param).Name() == "WhereParam" {
		slice := param.(WhereParam)
		where := make([]func(db *gorm.DB) *gorm.DB, 0, len(slice))
		for _, v := range slice {
			v := v
			where = append(where, func(db *gorm.DB) *gorm.DB {
				return db.Where(v.Sql, v.Value)
			})
		}
		return where
	}

	fields := structs.New(param).Fields()
	where := make([]func(db *gorm.DB) *gorm.DB, 0, len(fields))
	for _, field := range fields {
		if field.IsZero() {
			continue
		}

		if field.IsEmbedded() {
			embed := TransformWhereParamToGorm(field.Value())
			where = append(where, embed...)
			continue
		}

		sql := field.Tag("rdb")
		value := field.Value()

		where = append(where, func(db *gorm.DB) *gorm.DB {
			return db.Where(sql, value)
		})
	}

	return where
}

func UpdatedValue(before map[string]interface{}, after interface{}) map[string]interface{} {
	diff := make(map[string]interface{})
	tool := structs.New(after)
	tool.TagName = "gorm"
	for key, v := range tool.Map() {
		if before[key] != v {
			diff[key] = v
		}
	}
	return diff
}
