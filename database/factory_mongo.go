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

func NewReplicaSetMongo(cfg *ReplicaSetMongoConfig) (*mongo.Client, error) {
	cfg.setDefaultValue()

	wConcern := writeconcern.New(writeconcern.WMajority(), writeconcern.J(true), writeconcern.WTimeout(10*time.Second))
	rConcern := readconcern.Majority()

	opt := options.Client().
		SetServerSelectionTimeout(10 * time.Second).
		SetMaxPoolSize(cfg.MaxPoolSize).
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

type ReplicaSetMongoConfig struct {
	User        string
	Password    string
	Address     string // localhost:27018 or localhost:27017,localhost:27018,localhost:27019
	ReplicaSet  string
	MaxPoolSize uint64
}

func (c *ReplicaSetMongoConfig) setDefaultValue() {
	if c.MaxPoolSize <= 0 {
		const default_ = 8
		c.MaxPoolSize = default_
	}
}

// URI format:
// mongodb://root:1234@localhost:27018
// or
// mongodb://root:1234@localhost:27017,localhost:27018,localhost:27019
func (c *ReplicaSetMongoConfig) URI() string {
	return fmt.Sprintf("mongodb://%v:%v@%v", c.User, c.Password, c.Address)
}
