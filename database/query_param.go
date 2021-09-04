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
			sub := TransformQueryParamToGorm(field.Value())
			where = append(where, sub...)
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
