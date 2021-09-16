// +build integration

package database_test

import (
	"time"

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
	SqlID    string             `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" bson:"-"`
	MongoID  primitive.ObjectID `gorm:"-"                                bson:"_id"`
	Name     string             `gorm:"column:name;type:varchar(50)"     bson:"name"`
	NoTzTime time.Time          `gorm:"column:no_tz_time;type:timestamp" bson:"no_tz_time"`
	TzTime   time.Time          `gorm:"column:tz_time;type:timestamptz"  bson:"tz_time"`
	UpdateAt time.Time          `gorm:"column:update_at;type:timestamptz;autoUpdateTime"  bson:"update_at"`
}

func dbName() string {
	return "integration_test"
}

func mongoClient(cfg *database.ReplicaSetMongoConfig) *mongo.Client {
	if cfg == nil {
		cfg = &database.ReplicaSetMongoConfig{
			User:       "root",
			Password:   "1234",
			Address:    "localhost:27017,localhost:27019,localhost:27018",
			ReplicaSet: "dev-rs",
		}
	}

	client, err := database.NewReplicaSetMongo(cfg)
	if err != nil {
		panic(err)
	}
	return client
}

func mysqlGorm(cfg *database.RMDBConfig) *database.WrapperGorm {
	if cfg == nil {
		cfg = &database.RMDBConfig{
			User:     "root",
			Password: "1234",
			Host:     "localhost",
			Port:     "3306",
			Database: dbName(),
		}
	}

	db, err := database.NewGormMysql(cfg, true)
	if err != nil {
		panic(err)
	}
	return db
}

func pgGorm(cfg *database.RMDBConfig) *database.WrapperGorm {
	if cfg == nil {
		cfg = &database.RMDBConfig{
			User:     "root",
			Password: "1234",
			Host:     "localhost",
			Port:     "5432",
			Database: dbName(),
		}
	}

	db, err := database.NewGormPgsql(cfg, true)
	if err != nil {
		panic(err)
	}
	return db
}
