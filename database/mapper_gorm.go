package database

import (
	"github.com/fatih/structs"
	"gorm.io/gorm"
)

// TransformQueryParamToGorm param 只能傳入 struct, 會解析 struct tag key = rdb 的內容,
// 其 field 數值 若為 go 的初始值 會忽略此條件,
// 需要在 struct tag 填入 適當的 sql where 敘述.
//
// 若想查 xxx=false, 將型別定義為 pointer *bool;
// 若想查 xxx is null, 將型別定義為 string;
// 理論上全部數值可以都用 string 來表示.
//
// example:
//
// type query struct {
//   UserIDs     []int  rdb:"user_id in (?)"
//   IsAdmin     *bool  rdb:"is_admin = ?"
//   BlockedTime string rdb:"blocked_time is ?"
//   IsBlocked   string rdb:"is_blocked = ?"
// }
//
// query.BlockedTime="null"
// query.IsBlocked="false"
func TransformQueryParamToGorm(param interface{}) []func(db *gorm.DB) *gorm.DB {
	if param == nil {
		return nil
	}

	fields := structs.New(param).Fields()

	where := make([]func(db *gorm.DB) *gorm.DB, 0, len(fields))

	for _, field := range fields {
		if field.IsZero() {
			continue
		}

		if field.IsEmbedded() {
			embed := TransformQueryParamToGorm(field.Value())
			where = append(where, embed...)
			continue
		}

		key := field.Tag("rdb")
		value := field.Value()

		where = append(where, func(db *gorm.DB) *gorm.DB {
			return db.Where(key, value)
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
