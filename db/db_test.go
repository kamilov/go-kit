package db_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/kamilov/go-kit/db"
	"github.com/loghole/dbhook"
	_ "github.com/mattn/go-sqlite3"
)

type (
	hook struct {
		writer io.Writer
	}
	hookBefore struct{}
	hookAfter  struct{}
	hookError  struct{}
)

func (h *hook) Before(ctx context.Context, _ *dbhook.HookInput) (context.Context, error) {
	_, _ = io.WriteString(h.writer, "before")
	return ctx, nil
}

func (h *hookBefore) Before(ctx context.Context, _ *dbhook.HookInput) (context.Context, error) {
	return ctx, nil
}

func (h *hook) After(ctx context.Context, _ *dbhook.HookInput) (context.Context, error) {
	_, _ = io.WriteString(h.writer, "after")
	return ctx, nil
}

func (h *hookAfter) After(ctx context.Context, _ *dbhook.HookInput) (context.Context, error) {
	return ctx, nil
}

func (h *hook) Error(ctx context.Context, _ *dbhook.HookInput) (context.Context, error) {
	_, _ = io.WriteString(h.writer, "error")
	return ctx, nil
}

func (h *hookError) Error(ctx context.Context, _ *dbhook.HookInput) (context.Context, error) {
	return ctx, nil
}

func TestDB(t *testing.T) {
	_, err := db.New()
	if err == nil {
		t.Error("db.New() did not return an error")
	}

	_, err = db.New(db.WithConfigDSN("test://localhost:3306"))
	if err == nil {
		t.Error("db.New() did not return driver error")
	}

	var buf bytes.Buffer

	testDB, _ := db.New(
		db.WithConfigDSN("sqlite://:memory:"),
		db.WithHook(&hook{writer: &buf}),
		db.WithHookBefore(&hookBefore{}),
		db.WithHookAfter(&hookAfter{}),
		db.WithHookError(&hookError{}),
	)

	_, err = testDB.Exec("CREATE TABLE foo (Id INTEGER PRIMARY KEY)")
	if err != nil {
		t.Error(err)
	}

	t.Run("test transactional", func(t *testing.T) {
		t.Helper()

		t.Run("test transactional success", func(t *testing.T) {
			t.Helper()

			err = testDB.Transactional(func(tx *db.Tx) error {
				_, _ = tx.Exec("INSERT INTO foo (Id) VALUES (1)")
				_, _ = tx.Exec("INSERT INTO foo (Id) VALUES (2)")

				return nil
			})
			if err != nil {
				t.Errorf("transactional error: %v", err)
			}
		})

		t.Run("test transactional failure with return error", func(t *testing.T) {
			t.Helper()

			err = testDB.Transactional(func(tx *db.Tx) error {
				_, _ = tx.Exec("INSERT INTO foo (Id) VALUES (1)")
				_, _ = tx.Exec("INSERT INTO foo (Id) VALUES (2)")

				return errors.New("test error")
			})
			if err == nil {
				t.Error("transactional did not return an error")
			}
		})

		t.Run("test transactional failure with panic", func(t *testing.T) {
			t.Helper()

			defer func() {
				if r := recover(); r == nil {
					t.Error("transactional did not panic")
				}
			}()

			err = testDB.Transactional(func(tx *db.Tx) error {
				_, _ = tx.Exec("INSERT INTO foo (Id) VALUES (1)")

				panic("test error")
			})
		})
	})

	t.Run("test transactional context", func(t *testing.T) {
		t.Helper()

		ctx := context.Background()

		t.Run("test transactional success", func(t *testing.T) {
			t.Helper()

			err = testDB.TransactionalTx(ctx, nil, func(tx *db.Tx) error {
				_, _ = tx.Exec("INSERT INTO foo (Id) VALUES (1)")
				_, _ = tx.Exec("INSERT INTO foo (Id) VALUES (2)")

				return nil
			})
			if err != nil {
				t.Errorf("transactional error: %v", err)
			}
		})

		t.Run("test transactional failure with return error", func(t *testing.T) {
			t.Helper()

			err = testDB.TransactionalTx(ctx, nil, func(tx *db.Tx) error {
				_, _ = tx.Exec("INSERT INTO foo (Id) VALUES (1)")
				_, _ = tx.Exec("INSERT INTO foo (Id) VALUES (2)")

				return errors.New("test error")
			})
			if err == nil {
				t.Error("transactional did not return an error")
			}
		})

		t.Run("test transactional failure with panic", func(t *testing.T) {
			t.Helper()

			defer func() {
				if r := recover(); r == nil {
					t.Error("transactional did not panic")
				}
			}()

			err = testDB.TransactionalTx(ctx, nil, func(tx *db.Tx) error {
				_, _ = tx.Exec("INSERT INTO foo (Id) VALUES (?)", 1)

				panic("test error")
			})
		})
	})

	if !strings.Contains(buf.String(), "before") ||
		!strings.Contains(buf.String(), "after") ||
		!strings.Contains(buf.String(), "error") {
		t.Error("error hooks")
	}
}
