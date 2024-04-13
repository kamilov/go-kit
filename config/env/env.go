package env

import (
	"context"
	"encoding"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"strconv"

	"github.com/kamilov/go-kit/config"
	kitReflection "github.com/kamilov/go-kit/utils/reflect"
)

type (
	LookupFunc func(string) (string, bool)
	contextKey int
)

const (
	TagName contextKey = iota
	Prefix
	defaultTagName = "env"
)

//nolint:gochecknoglobals // used to redeclare default lookup function
var LookupEnv LookupFunc

//nolint:gochecknoinits // used for automatic adding read config function
func init() {
	LookupEnv = os.LookupEnv
	config.RegisterReader(reader)
}

func reader(ctx context.Context, data any) error {
	if err := kitReflection.ValidateStructPointer(data); err != nil {
		return err
	}

	tagName := getTagName(ctx)
	prefix := getPrefix(ctx)

	rv := reflect.Indirect(reflect.ValueOf(data))
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)

		if !field.CanSet() {
			continue
		}

		fieldType := rt.Field(i)

		if fieldType.Anonymous {
			field = reflect.Indirect(field)

			if field.Kind() == reflect.Struct {
				if err := reader(ctx, field.Addr().Interface()); err != nil {
					return err
				}
			}

			continue
		}

		name := fieldType.Tag.Get(tagName)
		if name == "-" {
			continue
		}

		name = prefix + name

		if value, ok := LookupEnv(name); ok {
			if err := set(field, value); err != nil {
				return err
			}
		}
	}

	return nil
}

func set(rv reflect.Value, value string) error {
	rv = reflect.Indirect(rv)
	rt := rv.Type()

	if !rv.CanAddr() {
		return errors.New("cannot assign address")
	}

	v := rv.Addr().Interface()
	if impl, ok := v.(encoding.TextUnmarshaler); ok {
		return impl.UnmarshalText([]byte(value))
	}
	if impl, ok := v.(encoding.BinaryUnmarshaler); ok {
		return impl.UnmarshalBinary([]byte(value))
	}

	//nolint:exhaustive // process only the listed types
	switch rt.Kind() {
	case reflect.String:
		rv.SetString(value)
	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		rv.SetBool(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		rv.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		rv.SetUint(val)
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}

		rv.SetFloat(val)
	case reflect.Slice:
		if rt.Elem().Kind() == reflect.Uint8 {
			sl := reflect.ValueOf([]byte(value))
			rv.Set(sl)
			return nil
		}
		fallthrough

	default:
		return json.Unmarshal([]byte(value), rv.Addr().Interface())
	}

	return nil
}

func getTagName(ctx context.Context) string {
	if value, ok := ctx.Value(TagName).(string); ok {
		return value
	}

	return defaultTagName
}

func getPrefix(ctx context.Context) string {
	if value, ok := ctx.Value(Prefix).(string); ok {
		return value
	}

	return ""
}
