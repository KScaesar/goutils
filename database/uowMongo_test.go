// +build integration

package database_test

import (
	"context"
	"testing"

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

	repo := bookMongoRepo{col: client.Database("golang_integration_test").Collection("testing_book")}
	uowFactory := database.NewUowMongoFactory(client)

	uow, err := uowFactory.CreateUow()
	assert.NoError(t, err)

	fn := func(txCtx context.Context) error {
		book1 := &DomainBook{Name: "ddd_is_good"}
		if err := repo.createBook(txCtx, book1); err != nil {
			return err
		}

		book1.Name = "tdd_is_good"
		if err := repo.updateBook(txCtx, book1); err != nil {
			return err
		}

		return nil
	}

	// enable tx
	uowErr := uow.AutoStart(nil, fn)
	assert.NoError(t, uowErr)

	// not enable transaction
	// _ = uow
	// repoErr := fn(nil)
	// assert.NoError(t, repoErr)
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
	filter := mongoFilter{}
	_, err := repo.col.UpdateOne(ctx, filter.ID(book.MongoID), bson.M{"$set": book})
	return err
}

type mongoFilter struct{}

func (f mongoFilter) ID(id interface{}) bson.D {
	return bson.D{{"_id", id}}
}
