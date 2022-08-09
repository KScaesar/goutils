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

	"github.com/KScaesar/goutils/database"
	"github.com/KScaesar/goutils/xLog"
)

func Test_txMongo_AutoComplete(t *testing.T) {
	client := mongoClient(nil)

	collectionName := "testing_books"
	repo := bookMongoRepo{
		col: client.
			Database(dbName).
			Collection(collectionName),
	}
	txFactory := database.NewMongoTxFactory(client, nil)

	tx := txFactory.CreateTx(nil)

	id := primitive.NewObjectID()
	fn := func(name string) func(txCtx context.Context) error {
		return func(txCtx context.Context) error {
			now := time.Now()
			book := &DomainBook{
				MongoID:  id,
				Name:     "python" + "#" + name,
				NoTzTime: now,
				TzTime:   now,
				// NullTime: goutils.Time(now),
				// AutoTIme: goutils.Time(now),
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

	// err = fn("disable tx")(nil)
	// assert.NoError(t, err, "disable tx")

	book, err := repo.getBook(nil, id)
	assert.NoError(t, err)
	xLog.Info().Interface("book", book).Send()
}

type bookMongoRepo struct {
	col *mongo.Collection
}

func (repo *bookMongoRepo) createBook(ctx context.Context, book *DomainBook) error {
	_, err := repo.col.InsertOne(ctx, book)
	return err
}

func (repo *bookMongoRepo) updateBook(ctx context.Context, book *DomainBook) error {
	_, err := repo.col.UpdateOne(ctx, bson.M{"id": book.MongoID}, bson.M{"$set": book})
	return err
}

func (repo *bookMongoRepo) getBook(ctx context.Context, id primitive.ObjectID) (book DomainBook, err error) {
	return book, repo.col.FindOne(ctx, bson.M{"_id": id}).Decode(&book)
}
