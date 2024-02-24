// Package db for working with database
package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/loghole/dbhook"
)

// DB struct to work with database
type DB struct {
	*sql.DB
}

// ErrUndefinedConfig configration error
var ErrUndefinedConfig = errors.New("undefined config")

// New open connection to database and return struct
func New(opts ...Option) (*DB, error) {
	//nolint:exhaustruct // define known fields
	o := options{
		hooks: make([]Hook, 0),
	}

	for _, opt := range opts {
		opt.apply(&o)
	}

	if o.config == nil {
		return nil, ErrUndefinedConfig
	}

	hooks := dbhook.NewHooks(dbhook.WithHook(o.hooks...))

	sqlDB, err := sql.Open(o.config.driverName(), "")
	if err != nil {
		return nil, fmt.Errorf("can't find original driver: %w", err)
	}

	replacedDriverName := fmt.Sprintf("%s-with-hook", o.config.driverName())

	sql.Register(replacedDriverName, dbhook.Wrap(sqlDB.Driver(), hooks))

	sqlDB, err = sql.Open(replacedDriverName, o.config.DSN())
	if err != nil {
		return nil, fmt.Errorf("can't open sql db: %w", err)
	}

	return &DB{sqlDB}, nil
}

// Begin open transaction
func (db *DB) Begin() (*Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin trasaction error: %w", err)
	}

	return &Tx{tx}, nil
}

// BeginTx open transaction with context
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("begin trasaction error: %w", err)
	}

	return &Tx{tx}, nil
}

// Transactional open transaction and run callback with queries and post process checking error
func (db *DB) Transactional(fn func(*Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("transaction start error: %w", err)
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
		} else {
			if err = tx.Commit(); errors.Is(err, sql.ErrTxDone) {
				err = nil
			}
		}
	}()

	err = fn(tx)

	return err
}

// TransactionalTx open transaction with context and run callback with queries and post process checking error
func (db *DB) TransactionalTx(ctx context.Context, opts *sql.TxOptions, fn func(*Tx) error) (err error) {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("begin transaction error: %w", err)
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
		} else {
			if err = tx.Commit(); errors.Is(err, sql.ErrTxDone) {
				err = nil
			}
		}
	}()

	err = fn(tx)

	return err
}
