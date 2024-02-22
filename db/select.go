package db

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"strings"

	"github.com/kamilov/go-kit/utils"
)

const tagName = "db"

var (
	ErrPointer     = errors.New("must be a pointer")
	ErrPointerType = errors.New("a pointer was not expected")
	fieldRegex     = regexp.MustCompile(`([^A-Z_])([A-Z])`)
)

func (db *DB) Select(ctx context.Context, data any, query string, args ...any) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	rv := reflect.ValueOf(data)

	if rv.Kind() != reflect.Ptr {
		return ErrPointer
	}

	rv = utils.Indirect(rv.Elem())

	switch rv.Kind() {
	case reflect.Slice:
		rows, err := db.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}

		for rows.Next() {
			val := reflect.New(rv.Type().Elem())

			if rv.Type().Elem().Kind() == reflect.Struct {
				scanStruct(rows, val)
			} else {
				_ = rows.Scan(val.Interface())
			}

			rv = reflect.Append(rv, val.Elem())
		}

	case reflect.Struct:
		rows, _ := db.QueryContext(ctx, query, args...)

		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return err
			}
			return sql.ErrNoRows
		}

		scanStruct(rows, rv)

	default:
		return ErrPointerType
	}

	return nil
}

func scanStruct(rows *sql.Rows, rv reflect.Value) {
	rv = utils.Indirect(rv)

	fieldNameIndex := map[string]int{}

	for i := 0; i < rv.Type().NumField(); i++ {
		var name string

		f := rv.Type().Field(i)
		tag := f.Tag.Get(tagName)

		if tag != "" {
			name = tag
		} else {
			name = fieldMap(f.Name)
		}

		fieldNameIndex[name] = i
	}

	columns, _ := rows.Columns()
	fields := make([]any, len(columns))

	for i, column := range columns {
		if fi, ok := fieldNameIndex[column]; ok {
			fields[i] = rv.Field(fi).Addr().Interface()
		} else {
			fields[i] = &sql.NullString{}
		}
	}

	_ = rows.Scan(fields...)
}

func fieldMap(name string) string {
	return strings.ToLower(fieldRegex.ReplaceAllString(name, "${1}_${2}"))
}
