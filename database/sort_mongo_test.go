package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/Min-Feng/goutils"
)

func TestMongoSort(t *testing.T) {
	type FooBarSortOption struct {
		CreateAt goutils.SortKind `bson:"create_at" xSort:"asc"`
		UpdateAt goutils.SortKind `bson:"update_at"`
		DeleteAt goutils.SortKind `bson:"delete_at" xSort:"desc"`
	}

	opt := FooBarSortOption{
		CreateAt: goutils.SortDesc,
	}
	expectedMongoSort := bson.D{
		{"create_at", -1},
		{"delete_at", -1},
	}

	actual := MongoSort(&opt)

	assert.Equal(t, expectedMongoSort, actual)
}
