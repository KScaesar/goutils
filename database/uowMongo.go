package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/Min-Feng/goutils/errorY"
)

func NewUowMongoFactory(client *mongo.Client) UowFactory {
	wConcern := writeconcern.New(writeconcern.WMajority(), writeconcern.J(true), writeconcern.WTimeout(10*time.Second))
	rConcern := readconcern.Snapshot()
	rPref := readpref.Primary()

	opt := options.Session().
		SetCausalConsistency(true).
		SetDefaultWriteConcern(wConcern).
		SetDefaultReadConcern(rConcern).
		SetDefaultReadPreference(rPref)

	return &uowMongoFactory{
		client:            client,
		createdSessionOpt: opt,
	}
}

type uowMongoFactory struct {
	client            *mongo.Client
	createdSessionOpt *options.SessionOptions
}

func (f *uowMongoFactory) CreateUow() (Uow, error) {
	session, err := f.client.StartSession(f.createdSessionOpt)
	if err != nil {
		return nil, errorY.Wrap(errorY.ErrSystem, err.Error())
	}
	return &uowMongo{sess: session}, nil
}

type uowMongo struct {
	sess mongo.Session
}

func (uow *uowMongo) AutoStart(ctx context.Context, txFn func(txCtx context.Context) error) error {
	if ctx == nil {
		ctx = context.Background()
	}
	defer uow.sess.EndSession(ctx)

	mongoTxFn := func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, txFn(sessCtx)
	}

	// https://docs.mongodb.com/v4.4/core/transactions-in-applications/#std-label-txn-callback-api
	_, err := uow.sess.WithTransaction(ctx, mongoTxFn)
	if err != nil {
		if errorY.IsUndefinedError(err) {
			return errorY.Wrap(errorY.ErrSystem, err.Error())
		}
		return err
	}

	return nil
}

func (uow *uowMongo) ManualStart(ctx context.Context, txFn func(txCtx context.Context) error) error {
	if ctx == nil {
		ctx = context.Background()
	}

	err := uow.sess.StartTransaction()
	if err != nil {
		return errorY.Wrap(errorY.ErrSystem, err.Error())
	}

	sessCtx := mongo.NewSessionContext(ctx, uow.sess)
	return txFn(sessCtx)
}

func (uow *uowMongo) Commit() error {
	err := uow.sess.CommitTransaction(nil)
	if err != nil {
		return errorY.Wrap(errorY.ErrSystem, err.Error())
	}

	uow.sess.EndSession(nil)
	return nil
}

func (uow *uowMongo) Rollback() error {
	defer uow.sess.EndSession(nil)

	err := uow.sess.AbortTransaction(nil)
	if err != nil {
		return errorY.Wrap(errorY.ErrSystem, err.Error())
	}
	return nil
}
