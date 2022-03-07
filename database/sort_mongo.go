package database

import (
	"github.com/fatih/structs"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/KScaesar/goutils"
)

// MongoSort
// 存在 xSort tag, 欄位為空值時, 填入預設值;
// 不存在 xSort tag, 欄位為空值時, 就跳過不採用此欄位
func MongoSort(sortOption interface{}) bson.D {
	if sortOption == nil {
		return nil
	}

	fields := structs.New(sortOption).Fields()
	where := make(bson.D, 0, len(fields))

	for _, field := range fields {
		value, ok := field.Value().(goutils.SortKind)
		if !ok {
			continue
		}

		if value == goutils.SortNone {
			tagV := field.Tag("xSort")
			if tagV == "" {
				continue
			}
			value = goutils.SortKind(tagV)
		}

		var mongoSortValue int
		switch value {
		case goutils.SortAsc:
			mongoSortValue = 1
		case goutils.SortDesc:
			mongoSortValue = -1
		}

		where = append(where, bson.E{
			Key:   field.Tag("bson"),
			Value: mongoSortValue,
		})
	}

	return where
}
