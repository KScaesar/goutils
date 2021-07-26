package database

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func NewMongo(cfg *MongoConfig) (*mongo.Client, error) {
	wConcern := writeconcern.New(writeconcern.WMajority(), writeconcern.J(true), writeconcern.WTimeout(10*time.Second))
	rConcern := readconcern.Majority()

	opt := options.Client().
		SetServerSelectionTimeout(10 * time.Second).
		SetMaxPoolSize(cfg.MaxPoolSize_()).
		SetWriteConcern(wConcern).
		SetReadConcern(rConcern).
		SetReplicaSet(cfg.ReplicaSet).
		ApplyURI(cfg.URI())

	client, err := mongo.Connect(nil, opt)
	if err != nil {
		return nil, err
	}

	err = client.Ping(nil, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client, nil
}

type MongoConfig struct {
	User        string
	Password    string
	Address     string
	ReplicaSet  string
	MaxPoolSize uint64
}

// URI format:
// mongodb://root:1234@localhost:27018
// or
// mongodb://root:1234@localhost:27017,localhost:27018,localhost:27019
func (c *MongoConfig) URI() string {
	return fmt.Sprintf("mongodb://%v:%v@%v", c.User, c.Password, c.Address)
}

func (c *MongoConfig) MaxPoolSize_() uint64 {
	if c.MaxPoolSize <= 0 {
		const defaultSize = 8
		return defaultSize
	}
	return c.MaxPoolSize
}
