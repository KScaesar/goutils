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

func NewMongoTxFactory(client *mongo.Client) TxFactory {
	wConcern := writeconcern.New(writeconcern.WMajority(), writeconcern.J(true), writeconcern.WTimeout(10*time.Second))
	rConcern := readconcern.Snapshot()
	rPref := readpref.Primary()

	opt := options.Session().
		SetCausalConsistency(true).
		SetDefaultWriteConcern(wConcern).
		SetDefaultReadConcern(rConcern).
		SetDefaultReadPreference(rPref)

	return &mongoTxFactory{
		client:            client,
		createdSessionOpt: opt,
	}
}

type mongoTxFactory struct {
	client            *mongo.Client
	createdSessionOpt *options.SessionOptions
}

func (f *mongoTxFactory) CreateTx() (Transaction, error) {
	session, err := f.client.StartSession(f.createdSessionOpt)
	if err != nil {
		return nil, errorY.Wrap(errorY.ErrSystem, err.Error())
	}
	return &mongoTxAdapter{sess: session}, nil
}

type mongoTxAdapter struct {
	sess mongo.Session
}

func (adapter *mongoTxAdapter) AutoComplete(ctx context.Context, fn func(txCtx context.Context) error) error {
	if ctx == nil {
		ctx = context.Background()
	}
	defer adapter.sess.EndSession(ctx)

	if adapter.isMongoSession(ctx) {
		return fn(ctx)
	}

	mongoTxFn := func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, fn(sessCtx)
	}

	// https://docs.mongodb.com/v4.4/core/transactions-in-applications/#std-label-txn-callback-api
	_, err := adapter.sess.WithTransaction(ctx, mongoTxFn)
	if err != nil {
		if errorY.IsUndefinedError(err) {
			return errorY.Wrap(errorY.ErrSystem, err.Error())
		}
		return err
	}

	return nil
}

func (adapter *mongoTxAdapter) ManualComplete(
	ctx context.Context,
	fn func(txCtx context.Context) error,
) (
	commit func() error,
	rollback func() error,
	fnErr error,
) {
	if ctx == nil {
		ctx = context.Background()
	}

	if adapter.isMongoSession(ctx) {
		return adapter.doNothing, adapter.doNothing, fn(ctx)
	}

	err := adapter.sess.StartTransaction()
	if err != nil {
		return adapter.doNothing, adapter.doNothing, errorY.Wrap(errorY.ErrSystem, err.Error())
	}

	sessCtx := mongo.NewSessionContext(ctx, adapter.sess)
	return adapter.commit, adapter.rollback, fn(sessCtx)
}

func (adapter *mongoTxAdapter) doNothing() error { return nil }

func (adapter *mongoTxAdapter) commit() error {
	defer adapter.sess.EndSession(nil)

	err := adapter.sess.CommitTransaction(nil)
	if err != nil {
		return errorY.Wrap(errorY.ErrSystem, err.Error())
	}

	return nil
}

func (adapter *mongoTxAdapter) rollback() error {
	defer adapter.sess.EndSession(nil)

	err := adapter.sess.AbortTransaction(nil)
	if err != nil {
		return errorY.Wrap(errorY.ErrSystem, err.Error())
	}

	return nil
}

func (adapter *mongoTxAdapter) isMongoSession(ctx context.Context) bool {
	_, ok := ctx.(mongo.SessionContext)
	return ok
}
