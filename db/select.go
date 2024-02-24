package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/kamilov/go-kit/utils"
)

const tagName = "db"

var (
	// ErrPointer error
	ErrPointer = errors.New("must be a pointer")
	// ErrPointerType error
	ErrPointerType = errors.New("a pointer was not expected")
	fieldRegex     = regexp.MustCompile(`([^A-Z_])([A-Z])`)
)

// Select data and parse to struct
//
//nolint:cyclop // normal cyclomatic
func (db *DB) Select(ctx context.Context, data any, query string, args ...any) error {
	if ctx.Err() != nil {
		return ctx.Err() //nolint:wrapcheck // return context error
	}

	rv := reflect.ValueOf(data)

	if rv.Kind() != reflect.Ptr {
		return ErrPointer
	}

	rv = utils.Indirect(rv.Elem())

	//nolint:exhaustive // process only struct and slice
	switch rv.Kind() {
	case reflect.Slice:
		rows, err := db.QueryContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("query error: %w", err)
		}

		defer func() {
			_ = rows.Close()
		}()

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

		defer func() {
			_ = rows.Close()
		}()

		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return fmt.Errorf("parse row error: %w", err)
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
			fields[i] = &sql.NullString{} //nolint:exhaustruct // not use fields
		}
	}

	_ = rows.Scan(fields...)
}

func fieldMap(name string) string {
	return strings.ToLower(fieldRegex.ReplaceAllString(name, "${1}_${2}"))
}
