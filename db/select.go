package db

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"strings"

	kitReflect "github.com/kamilov/go-kit/utils/reflect"
)

const tagName = "db"

var (
	ErrPointerType        = errors.New("a pointer was not expected")
	camelCaseToSnakeRegex = regexp.MustCompile(`([^A-Z_])([A-Z])`)
)

func (db *DB) Select(ctx context.Context, data any, query string, args ...any) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if err := kitReflect.ValidatePointer(data); err != nil {
		return err
	}

	rp := reflect.ValueOf(data)
	rv := reflect.Indirect(rp)

	//nolint:exhaustive // process only struct and slice
	switch rv.Kind() {
	case reflect.Struct:
		rows, err := db.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}

		defer func() {
			_ = rows.Close()
		}()

		if !rows.Next() {
			if err = rows.Err(); err != nil {
				return err
			}

			return sql.ErrNoRows
		}

		scanStruct(rows, rv)

	case reflect.Slice:
		rows, err := db.QueryContext(ctx, query, args...)
		if err != nil {
			return err
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

		if rp.Elem().CanSet() {
			rp.Elem().Set(rv)
		}
	default:
		return ErrPointerType
	}

	return nil
}

func scanStruct(rows *sql.Rows, rv reflect.Value) {
	rv = reflect.Indirect(rv)
	fieldNameIndex := make(map[string]int)

	for i := 0; i < rv.Type().NumField(); i++ {
		var name string

		field := rv.Type().Field(i)
		tag := field.Tag.Get(tagName)

		switch tag {
		case "-":
			continue
		case "":
			name = fieldMap(field.Name)
		default:
			name = tag
		}

		fieldNameIndex[name] = i
	}

	columns, _ := rows.Columns()
	fields := make([]any, len(columns))

	for i, col := range columns {
		if index, ok := fieldNameIndex[col]; ok {
			fields[i] = rv.Field(index).Addr().Interface()
		} else {
			fields[i] = &sql.NullString{}
		}
	}

	_ = rows.Scan(fields...)
}

func fieldMap(name string) string {
	return strings.ToLower(camelCaseToSnakeRegex.ReplaceAllString(name, "${1}_${2}"))
}
