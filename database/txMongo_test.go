//go:build integration
// +build integration

package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Min-Feng/goutils"
	"github.com/Min-Feng/goutils/database"
)

func Test_txMongo_AutoComplete(t *testing.T) {
	client := mongoClient(nil)

	mongoBook := infraBook{}
	repo := bookMongoRepo{
		col: client.
			Database(dbName).
			Collection(mongoBook.CollectionName()),
	}
	txFactory := database.NewMongoTxFactory(client, nil)

	tx := txFactory.CreateTx(nil)

	fn := func(name string) func(txCtx context.Context) error {
		return func(txCtx context.Context) error {
			book := &DomainBook{
				SqlID:    "",
				MongoID:  primitive.ObjectID{},
				Name:     "python" + "#" + name,
				NoTzTime: time.Now(),
				TzTime:   time.Now(),
				UpdateAt: goutils.Time(time.Now()),
			}
			if err := repo.createBook(txCtx, book); err != nil {
				return err
			}

			book.Name = "golang" + "#" + name
			if err := repo.updateBook(txCtx, book); err != nil {
				return err
			}

			return nil
		}
	}

	err := tx.AutoComplete(fn("enable tx"))
	assert.NoError(t, err, "enable tx")

	err = fn("disable tx")(nil)
	assert.NoError(t, err, "disable tx")
}

type bookMongoRepo struct {
	col *mongo.Collection
}

func (repo *bookMongoRepo) createBook(ctx context.Context, book *DomainBook) error {
	book.MongoID = primitive.NewObjectID()
	_, err := repo.col.InsertOne(ctx, book)
	return err
}

func (repo *bookMongoRepo) updateBook(ctx context.Context, book *DomainBook) error {
	query := mongoQuery{}
	_, err := repo.col.UpdateOne(ctx, query.ID(book.MongoID), bson.M{"$set": book})
	return err
}

type mongoQuery struct{}

func (f mongoQuery) ID(id interface{}) bson.D {
	return bson.D{{"_id", id}}
}
