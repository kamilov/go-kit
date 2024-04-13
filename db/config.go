package db

import (
	"bytes"
	"net/url"
	"sort"
)

type driverName string

const (
	PGX        driverName = "pgx"
	Postgres   driverName = "postgres"
	Clickhouse driverName = "clickhouse"
	SQLite     driverName = "sqlite3"
)

type Config struct {
	Hostname string
	Username string
	Password string
	Database string
	Driver   driverName
	Params   map[string]string
}

func (c *Config) DSN() string {
	switch c.Driver {
	case PGX, Postgres:
		return c.postgresConnectionString()
	case Clickhouse:
		return c.clickhouseConnectionString()
	case SQLite:
		return c.sqliteConnectionString()
	default:
		panic("unknown Driver name")
	}
}

func (c *Config) driverName() string {
	return string(c.Driver)
}

func (c *Config) postgresConnectionString() string {
	var dsn url.URL

	dsn.Scheme = "postgres"

	if c.Password != "" {
		dsn.User = url.UserPassword(c.Username, c.Password)
	} else {
		dsn.User = url.User(c.Username)
	}

	dsn.Host = c.Hostname
	dsn.RawQuery = encodeParams(c.Params)

	return dsn.JoinPath(c.Database).String()
}

func (c *Config) clickhouseConnectionString() string {
	params := c.Params

	if params == nil {
		params = map[string]string{}
	}

	if c.Username != "" {
		params["username"] = c.Username
	}

	if c.Password != "" {
		params["password"] = c.Password
	}

	var dsn url.URL

	dsn.Scheme = "clickhouse"
	dsn.Host = c.Hostname
	dsn.RawQuery = encodeParams(params)

	return dsn.JoinPath(c.Database).String()
}

func (c *Config) sqliteConnectionString() string {
	if len(c.Params) == 0 {
		return c.Database
	}

	return c.Database + "?" + encodeParams(c.Params)
}

func encodeParams(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}

	keys := make([]string, 0, len(params))

	for key := range params {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	var buf bytes.Buffer

	for _, key := range keys {
		value := params[key]

		if buf.Len() > 0 {
			_ = buf.WriteByte('&')
		}

		_, _ = buf.WriteString(key)
		_ = buf.WriteByte('=')
		_, _ = buf.WriteString(value)
	}

	return buf.String()
}
