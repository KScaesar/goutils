// +build integration

package database_test

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Min-Feng/goutils/database"
)

type infraBook struct {
	// DomainBook 必須設為 public
	// 才能夠被 gorm.Migrator().CreateTable 感知
	// 因為 reflect 無法對 unexported field 進行處理
	DomainBook
}

func (s *infraBook) TableName() string {
	return "testing_books"
}

func (s *infraBook) collectionName() string {
	return "testing_books"
}

type DomainBook struct {
	SqlID   int                `gorm:"column:id"                    bson:"-"`
	MongoID primitive.ObjectID `gorm:"-"                            bson:"_id"`
	Name    string             `gorm:"column:name;type:varchar(50)" bson:"name"`
}

type testFixture struct{}

func (f testFixture) databaseName() string {
	return "golang_integration_test"
}

func (f testFixture) mongoConnectConfig() database.MongoConfig {
	return database.MongoConfig{
		User:       "root",
		Password:   "1234",
		Address:    "localhost:27017,localhost:27019,localhost:27018",
		ReplicaSet: "dev-rs",
	}
}

func (f testFixture) mongoClient(cfg database.MongoConfig) *mongo.Client {
	client, err := database.NewMongo(&cfg)
	if err != nil {
		panic(err)
	}
	return client
}

func (f testFixture) mysqlConnectConfig() database.RMDBConfig {
	return database.RMDBConfig{
		User:        "root",
		Password:    "1234",
		HostAndPort: "localhost:3306",
		Database:    f.databaseName(),
	}
}

func (f testFixture) mysqlGorm(cfg database.RMDBConfig) *database.WrapperGorm {
	db, err := database.NewGormMysql(&cfg)
	if err != nil {
		panic(err)
	}
	return db
}
