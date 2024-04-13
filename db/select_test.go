package db_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/kamilov/go-kit/db"
	_ "github.com/mattn/go-sqlite3"
)

func TestSelect(t *testing.T) {
	testDB, _ := db.New(
		db.WithConfigDSN("sqlite://:memory:"),
	)

	_, err := testDB.Exec("CREATE TABLE foo (id INTEGER PRIMARY KEY)")
	if err != nil {
		t.Fatal(err)
	}

	for i := 1; i <= 10; i++ {
		_, _ = testDB.Exec("INSERT INTO foo (id) VALUES (?)", i)
	}

	ctx := context.Background()

	t.Run("test select slice", func(t *testing.T) {
		t.Helper()

		t.Run("int", func(t *testing.T) {
			t.Helper()

			ids := make([]int, 0)
			err = testDB.Select(ctx, &ids, "SELECT id FROM foo")
			if err != nil {
				t.Error(err)
			} else if len(ids) != 10 {
				t.Errorf("got %d ids, want 10", len(ids))
			}
		})

		t.Run("string", func(t *testing.T) {
			t.Helper()

			ids := make([]string, 0)
			err = testDB.Select(ctx, &ids, "SELECT id FROM foo")
			if err != nil {
				t.Error(err)
			} else if len(ids) != 10 {
				t.Errorf("got %d ids, want 10", len(ids))
			}
		})

		t.Run("struct", func(t *testing.T) {
			t.Helper()

			ids := make([]struct {
				Value string `db:"id"`
			}, 0, 5)
			err = testDB.Select(ctx, &ids, "SELECT id FROM foo LIMIT 5")
			switch {
			case err != nil:
				t.Error(err)
			case len(ids) != 5:
				t.Errorf("got %d ids length, want 5", len(ids))
			case ids[0].Value != "1":
				t.Errorf("got %q, want %q", ids[0].Value, "1")
			}
		})
	})

	t.Run("test select struct", func(t *testing.T) {
		t.Helper()

		var row struct {
			ID uint64
		}

		err = testDB.Select(ctx, &row, "SELECT id FROM foo LIMIT 1")
		if err != nil {
			t.Error(err)
		} else if row.ID != 1 {
			t.Errorf("got %d, want 1", row.ID)
		}
	})

	t.Run("test select context error", func(t *testing.T) {
		t.Helper()

		str := ""
		err = testDB.Select(ctx, str, "SELECT id FROM foo LIMIT 1")
		if err == nil {
			t.Error("expected error pointer error")
		}

		err = testDB.Select(ctx, &str, "SELECT id FROM foo LIMIT 1")
		if err == nil {
			t.Error("expected error pointer type error")
		}
	})

	t.Run("test select struct not found", func(t *testing.T) {
		t.Helper()

		var row struct {
			ID uint64
		}

		err = testDB.Select(ctx, &row, "SELECT id FROM foo WHERE id = 11")
		if !errors.Is(err, sql.ErrNoRows) {
			t.Error("expected sql.ErrNoRows")
		}
	})
}
