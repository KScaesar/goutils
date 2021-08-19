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

func Test_uowMongo_AutoStart(t *testing.T) {
	fixture := testFixture{}
	config := fixture.mongoConnectConfig()
	client := fixture.mongoClient(config)

	mongoBook := infraBook{}
	repo := bookMongoRepo{
		col: client.
			Database(fixture.dbName()).
			Collection(mongoBook.collectionName()),
	}
	uowFactory := database.NewUowMongoFactory(client)

	uow, err := uowFactory.CreateUow()
	assert.NoError(t, err)

	createFn := func(name string) func(txCtx context.Context) error {
		return func(txCtx context.Context) error {
			book := &DomainBook{
				Name:     "ddd_is_good" + "#" + name,
				NoTzTime: time.Now(),
				TzTime:   time.Now(),
			}
			if err := repo.createBook(txCtx, book); err != nil {
				return err
			}

			book.Name = "tdd_is_good" + "#" + name
			if err := repo.updateBook(txCtx, book); err != nil {
				return err
			}

			return nil
		}
	}

	txFn := createFn("tx")
	uowErr := uow.AutoStart(nil, txFn)
	assert.NoError(t, uowErr, "enable tx")

	noTxFn := createFn("noTx")
	fnErr := noTxFn(nil)
	assert.NoError(t, fnErr, "not enable transaction")
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
