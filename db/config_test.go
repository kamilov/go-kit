package db_test

import (
	"testing"

	"github.com/kamilov/go-kit/db"
)

func TestConfig(t *testing.T) {
	t.Run("test postgres config", func(t *testing.T) {
		t.Helper()

		cfg := &db.Config{
			Hostname: "localhost",
			Username: "postgres",
			Password: "postgres",
			Database: "postgres",
			Driver:   db.Postgres,
		}

		if cfg.DSN() != "postgres://postgres:postgres@localhost/postgres" {
			t.Error("invalid Postgres DSN")
		}

		cfg = &db.Config{
			Hostname: "localhost",
			Username: "postgres",
			Database: "postgres",
			Driver:   db.PGX,
			Params:   map[string]string{"sslmode": "disable", "charset": "utf8"},
		}

		if cfg.DSN() != "postgres://postgres@localhost/postgres?charset=utf8&sslmode=disable" {
			t.Errorf("invalid PGX DSN: %v", cfg.DSN())
		}
	})

	t.Run("test clickhouse config", func(t *testing.T) {
		t.Helper()

		cfg := &db.Config{
			Hostname: "localhost",
			Username: "clickhouse",
			Password: "clickhouse",
			Database: "clickhouse",
			Driver:   db.Clickhouse,
		}

		if cfg.DSN() != "clickhouse://localhost/clickhouse?password=clickhouse&username=clickhouse" {
			t.Errorf("invalid ClickHouse DSN: %v", cfg.DSN())
		}
	})

	t.Run("test sqlite config", func(t *testing.T) {
		t.Helper()

		cfg := &db.Config{
			Database: ":memory:",
			Driver:   db.SQLite,
		}

		if cfg.DSN() != ":memory:" {
			t.Errorf("invalid SQLite DSN: %v", cfg.DSN())
		}

		cfg = &db.Config{
			Database: ":memory:",
			Driver:   db.SQLite,
			Params:   map[string]string{"charset": "utf8"},
		}

		if cfg.DSN() != ":memory:?charset=utf8" {
			t.Errorf("invalid SQLite DSN: %v", cfg.DSN())
		}
	})

	t.Run("test unknown driver", func(t *testing.T) {
		t.Helper()

		cfg := &db.Config{}

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("driver detection error: %v", cfg.Driver)
			}
		}()

		_ = cfg.DSN()
	})
}
