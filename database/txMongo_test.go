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

	"github.com/Min-Feng/goutils/database"
)

func Test_txMongo_AutoComplete(t *testing.T) {
	client := mongoClient(nil)

	mongoBook := infraBook{}
	repo := bookMongoRepo{
		col: client.
			Database(dbName()).
			Collection(mongoBook.collectionName()),
	}
	txFactory := database.NewMongoTxFactory(client)

	tx, err := txFactory.CreateTx(nil)
	assert.NoError(t, err)

	fn := func(name string) func(txCtx context.Context) error {
		return func(txCtx context.Context) error {
			book := &DomainBook{
				Name:     "python" + "#" + name,
				NoTzTime: time.Now(),
				TzTime:   time.Now(),
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

	err = tx.AutoComplete(fn("tx"))
	assert.NoError(t, err, "enable tx")

	err = fn("noTx")(nil)
	assert.NoError(t, err, "not enable tx")
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
