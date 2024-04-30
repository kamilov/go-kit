package text

import (
	"bytes"
	"context"
	"encoding"
	"io"
	"reflect"
	"strconv"
	"sync"
)

type unmarshalFunc = func([]byte, reflect.Value) error

var (
	unmarshalType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem() //nolint:gochecknoglobals // use global
	unmarshalers  sync.Map                                                  //nolint:gochecknoglobals // use global
)

func decoder(_ context.Context, reader io.Reader, data any) error {
	body, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	body = bytes.TrimSpace(body)

	switch data := data.(type) {
	case *string:
		*data = string(body)

	case *[]byte:
		*data = body

	case *int:
		return decodeInt(body, data, numDefaultSize)

	case *int8:
		return decodeInt(body, data, num8BitSize)

	case *int16:
		return decodeInt(body, data, num16BitSize)

	case *int32:
		return decodeInt(body, data, num32BitSize)

	case *int64:
		return decodeInt(body, data, num64BitSize)

	case *uint:
		return decodeUint(body, data, numDefaultSize)

	case *uint8:
		return decodeUint(body, data, num8BitSize)

	case *uint16:
		return decodeUint(body, data, num16BitSize)

	case *uint32:
		return decodeUint(body, data, num32BitSize)

	case *uint64:
		return decodeUint(body, data, num64BitSize)

	case *float32:
		return decodeFloat(body, data, num32BitSize)

	case *float64:
		return decodeFloat(body, data, num64BitSize)

	case *bool:
		*data, err = strconv.ParseBool(string(body))

	default:
		return unmarshal(body, data)
	}

	return err
}

func decodeInt[T int | int8 | int16 | int32 | int64](b []byte, data *T, bits int) error {
	i, err := strconv.ParseInt(string(b), numBase, bits)
	*data = T(i)

	return err
}

func decodeUint[T uint | uint8 | uint16 | uint32 | uint64](b []byte, data *T, bits int) error {
	i, err := strconv.ParseUint(string(b), numBase, bits)
	*data = T(i)

	return err
}

func decodeFloat[T float32 | float64](b []byte, data *T, bits int) error {
	f, err := strconv.ParseFloat(string(b), bits)
	*data = T(f)

	return err
}

func unmarshal(b []byte, data any) error {
	rv := reflect.ValueOf(data)
	rt := rv.Type()

	if dec, ok := unmarshalers.Load(rt); ok {
		return dec.(unmarshalFunc)(b, rv)
	}

	dec, err := newUnmarshaler(rt)
	if err != nil {
		return err
	}

	unmarshalers.Store(rt, dec)

	return dec(b, rv)
}

func newUnmarshaler(rt reflect.Type) (unmarshalFunc, error) {
	if rt.Implements(unmarshalType) {
		isPointer := rt.Kind() == reflect.Pointer
		rt = rt.Elem()

		return func(bytes []byte, rv reflect.Value) error {
			if isPointer && rv.IsNil() {
				rv.Set(reflect.New(rt))
			}

			return rv.Interface().(encoding.TextUnmarshaler).UnmarshalText(bytes)
		}, nil
	}

	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
		dec, err := newUnmarshaler(rt)
		if err != nil {
			return nil, err
		}

		return func(bytes []byte, rv reflect.Value) error {
			if rv.IsNil() {
				rv.Set(reflect.New(rt))
			}

			return dec(bytes, rv.Elem())
		}, nil
	}

	return nil, ErrUnsupportedType
}
