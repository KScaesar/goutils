package database

import (
	"github.com/fatih/structs"
	"gorm.io/gorm"
)

func TransformQueryParamToGorm(param interface{}) []func(db *gorm.DB) *gorm.DB {
	if param == nil {
		return nil
	}

	tool := structs.New(param)
	fields := tool.Fields()

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

		key := field.Tag("xQuery")
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
