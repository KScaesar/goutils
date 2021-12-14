package database

import (
	"github.com/fatih/structs"
	"gorm.io/gorm"

	"github.com/Min-Feng/goutils"
)

// GormFilter filter 可傳入 struct or goutils.FilterOption,
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
func GormFilter(filter interface{}) []func(db *gorm.DB) *gorm.DB {
	if filter == nil {
		return nil
	}

	{
		elements, ok := filter.(goutils.FilterOption)
		if ok {
			where := make([]func(db *gorm.DB) *gorm.DB, 0, len(elements))
			for _, e := range elements {
				element := e
				where = append(where, func(db *gorm.DB) *gorm.DB {
					return db.Where(element.Key, element.Value)
				})
			}
			return where
		}
	}

	{
		fields := structs.New(filter).Fields()
		where := make([]func(db *gorm.DB) *gorm.DB, 0, len(fields))

		for _, field := range fields {
			if field.IsZero() {
				continue
			}

			value := field.Value()
			if field.IsEmbedded() {
				embed := GormFilter(value)
				where = append(where, embed...)
				continue
			}

			sql := field.Tag("rdb")
			where = append(where, func(db *gorm.DB) *gorm.DB {
				return db.Where(sql, value)
			})
		}

		return where
	}
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
