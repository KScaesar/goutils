package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/mongo/driver"

	"github.com/Min-Feng/goutils/errors"
)

type MongoTxConfigs struct {
	SessionTimeout time.Duration
}

func (c *MongoTxConfigs) setDefaultValue() {
	if c.SessionTimeout <= 0 {
		const default_ = 10 * time.Second
		c.SessionTimeout = default_
	}
}

func NewMongoTxFactory(client *mongo.Client, cfg *MongoTxConfigs) TxFactory {
	if cfg == nil {
		cfg = new(MongoTxConfigs)
	}
	cfg.setDefaultValue()

	opt := NewMongoSessionOptByTransaction(cfg.SessionTimeout)

	return &mongoTxFactory{
		client: client,
		opt:    opt,
	}
}

func NewMongoSessionOptByTransaction(timeout time.Duration) *options.SessionOptions {
	wConcern := writeconcern.New(
		writeconcern.WMajority(),
		writeconcern.J(true),
		writeconcern.WTimeout(timeout),
	)
	rConcern := readconcern.Snapshot()
	rPref := readpref.Primary()

	return options.Session().
		SetCausalConsistency(true).
		SetDefaultWriteConcern(wConcern).
		SetDefaultReadConcern(rConcern).
		SetDefaultReadPreference(rPref)
}

type mongoTxFactory struct {
	client *mongo.Client
	opt    *options.SessionOptions
}

func (f *mongoTxFactory) CreateTx(ctx context.Context) Transaction {
	if ctx == nil {
		ctx = context.Background()
	}
	return &mongoTxAdapter{
		client: f.client,
		opt:    f.opt,
		ctx:    ctx,
	}
}

type mongoTxAdapter struct {
	client *mongo.Client
	opt    *options.SessionOptions

	sess mongo.Session
	ctx  context.Context
}

func (adapter *mongoTxAdapter) AutoComplete(fn func(txCtx context.Context) error) error {
	ctx := adapter.ctx

	if adapter.ExistMongoSessionInsideContext(ctx) {
		return fn(ctx)
	}

	if err := adapter.createSession(); err != nil {
		return err
	}
	defer adapter.sess.EndSession(ctx)

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

func (adapter *mongoTxAdapter) createSession() error {
	session, err := adapter.client.StartSession(adapter.opt)
	if err != nil {
		return errors.Wrap(errors.ErrSystem, err.Error())
	}

	adapter.sess = session
	return nil
}

func (adapter *mongoTxAdapter) ManualComplete(fn func(txCtx context.Context) error) (
	commit func() error,
	rollback func() error,
	fnErr error,
) {
	ctx := adapter.ctx

	if adapter.ExistMongoSessionInsideContext(ctx) {
		return adapter.doNothing, adapter.doNothing, fn(ctx)
	}

	if err := adapter.createSession(); err != nil {
		return adapter.doNothing, adapter.doNothing, err
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

func (adapter *mongoTxAdapter) ExistMongoSessionInsideContext(ctx context.Context) bool {
	_, ok := ctx.(mongo.SessionContext)
	return ok
}

func TranslateMongoError(err error) error {
	if err == nil {
		return nil
	}

	var cmdErr *mongo.CommandError

	switch {
	// https://docs.mongodb.com/v4.4/core/transactions-in-applications/#std-label-txn-callback-api
	case errors.As(err, &cmdErr) && cmdErr.HasErrorLabel(driver.TransientTransactionError):
		return err

	case errors.Is(err, mongo.ErrNoDocuments):
		return errors.Wrap(errors.ErrNotFound, err.Error())

	default:
		return errors.Wrap(errors.ErrSystem, err.Error())
	}
}
