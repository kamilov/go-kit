package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/loghole/dbhook"
)

type DB struct {
	*sql.DB
}

func New(opts ...Option) (*DB, error) {
	o := options{
		hooks: make([]Hook, 0),
	}

	for _, opt := range opts {
		opt.apply(&o)
	}

	if o.config == nil {
		return nil, errors.New("undefined config")
	}

	hooks := dbhook.NewHooks(dbhook.WithHook(o.hooks...))
	sqlDB, err := sql.Open(o.config.DriverName(), "")
	if err != nil {
		return nil, fmt.Errorf("can't find origional driver: %w", err)
	}

	replacedDriverName := fmt.Sprintf("%s-with-hook", o.config.DriverName())

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

func (db *DB) Transactional(fn func(*Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
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

func (db *DB) TransactionalContext(ctx context.Context, opts *sql.TxOptions, fn func(*Tx) error) (err error) {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
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
