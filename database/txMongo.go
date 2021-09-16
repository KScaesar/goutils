package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/Min-Feng/goutils/errors"
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

func (f *mongoTxFactory) CreateTx(ctx context.Context) (Transaction, error) {
	session, err := f.client.StartSession(f.createdSessionOpt)
	if err != nil {
		return nil, errors.Wrap(errors.ErrSystem, err.Error())
	}

	if ctx == nil {
		ctx = context.Background()
	}

	return &mongoTxAdapter{
		sess: session,
		ctx:  ctx,
	}, nil
}

type mongoTxAdapter struct {
	sess mongo.Session
	ctx  context.Context
}

func (adapter *mongoTxAdapter) AutoComplete(fn func(txCtx context.Context) error) error {
	ctx := adapter.ctx
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
		if errors.IsUndefinedError(err) {
			return errors.Wrap(errors.ErrSystem, err.Error())
		}
		return err
	}

	return nil
}

func (adapter *mongoTxAdapter) ManualComplete(fn func(txCtx context.Context) error) (
	commit func() error,
	rollback func() error,
	fnErr error,
) {
	ctx := adapter.ctx

	if adapter.isMongoSession(ctx) {
		adapter.sess.EndSession(ctx)
		return adapter.doNothing, adapter.doNothing, fn(ctx)
	}

	err := adapter.sess.StartTransaction()
	if err != nil {
		adapter.sess.EndSession(ctx)
		return adapter.doNothing, adapter.doNothing, errors.Wrap(errors.ErrSystem, err.Error())
	}

	sessCtx := mongo.NewSessionContext(ctx, adapter.sess)
	return adapter.commit, adapter.rollback, fn(sessCtx)
}

func (adapter *mongoTxAdapter) doNothing() error { return nil }

func (adapter *mongoTxAdapter) commit() error {
	ctx := adapter.ctx
	defer adapter.sess.EndSession(ctx)

	err := adapter.sess.CommitTransaction(ctx)
	if err != nil {
		return errors.Wrap(errors.ErrSystem, err.Error())
	}

	return nil
}

func (adapter *mongoTxAdapter) rollback() error {
	ctx := adapter.ctx
	defer adapter.sess.EndSession(ctx)

	err := adapter.sess.AbortTransaction(ctx)
	if err != nil {
		return errors.Wrap(errors.ErrSystem, err.Error())
	}

	return nil
}

func (adapter *mongoTxAdapter) isMongoSession(ctx context.Context) bool {
	_, ok := ctx.(mongo.SessionContext)
	return ok
}
