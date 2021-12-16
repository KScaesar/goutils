package database

import (
	"reflect"

	"github.com/fatih/structs"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/Min-Feng/goutils"
)

func MongoFilter(filter interface{}) bson.D {
	if filter == nil {
		return nil
	}

	{
		elements, ok := filter.(goutils.FilterOption)
		if ok {
			where := make(bson.D, 0, len(elements))
			for _, e := range elements {
				where = append(where, bson.E{
					Key:   e.Key,
					Value: e.Value,
				})
			}
			return where
		}
	}

	{
		fields := structs.New(filter).Fields()
		where := make(bson.D, 0, len(fields))

		for _, field := range fields {
			if field.IsZero() {
				continue
			}

			value := field.Value()
			if field.IsEmbedded() {
				embed := MongoFilter(value)
				where = append(where, embed...)
				continue
			}

			rv := reflect.ValueOf(value)
			key := field.Tag("bson")

			if rv.Kind() == reflect.Slice {
				where = append(where, bson.E{
					Key:   key,
					Value: bson.M{"$in": value},
				})
				continue
			}

			where = append(where, bson.E{
				Key:   key,
				Value: value,
			})
		}

		return where
	}
}
