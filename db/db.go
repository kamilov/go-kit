package db

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"

	"github.com/loghole/dbhook"
)

type DB struct {
	*sql.DB
}

const driverNameRandomSize = 5

var ErrUndefinedConfig = errors.New("configuration is not specified")

func New(opts ...Option) (*DB, error) {
	o := &options{
		hookOptions: make([]dbhook.HookOption, 0),
	}

	for _, opt := range opts {
		opt.apply(o)
	}

	if o.config == nil {
		return nil, ErrUndefinedConfig
	}

	hooks := dbhook.NewHooks(o.hookOptions...)
	sqlDB, err := sql.Open(o.config.driverName(), "")
	if err != nil {
		return nil, fmt.Errorf("can't find original driver: %w", err)
	}

	randomBytes := make([]byte, driverNameRandomSize)
	if _, err = rand.Read(randomBytes); err != nil {
		return nil, fmt.Errorf("can't read random bytes: %w", err)
	}
	replacedDriverName := fmt.Sprintf("%s-with-hooks-%s", o.config.driverName(), randomBytes)

	sql.Register(replacedDriverName, dbhook.Wrap(sqlDB.Driver(), hooks))

	sqlDB, err = sql.Open(replacedDriverName, o.config.DSN())
	if err != nil {
		return nil, fmt.Errorf("can't open sql db: %w", err)
	}

	return &DB{sqlDB}, nil
}

func (db *DB) Begin() (*Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}

	return &Tx{tx}, nil
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &Tx{tx}, nil
}

func (db *DB) Transactional(fn func(tx *Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		//nolint:nestif // no check
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			if err2 := tx.Rollback(); err2 != nil {
				if errors.Is(err2, sql.ErrTxDone) {
					return
				}
				err = err2
			}
		} else if err = tx.Commit(); errors.Is(err, sql.ErrTxDone) {
			err = nil
		}
	}()

	err = fn(tx)

	return err
}

func (db *DB) TransactionalTx(ctx context.Context, opts *sql.TxOptions, fn func(tx *Tx) error) (err error) {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		//nolint:nestif // no check
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			if err2 := tx.Rollback(); err2 != nil {
				if errors.Is(err2, sql.ErrTxDone) {
					return
				}
				err = err2
			}
		} else if err = tx.Commit(); errors.Is(err, sql.ErrTxDone) {
			err = nil
		}
	}()

	err = fn(tx)

	return err
}
