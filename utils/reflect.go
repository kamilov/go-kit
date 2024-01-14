package utils

import (
	"encoding"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
)

// Setter interface describing types can determine its self value
type Setter interface {
	Set(value string) error
}

// ErrCanAddr error for unaddressable cases
var ErrCanAddr = errors.New("the value is unaddressable")

// Indirect unpoint reflect value
func Indirect(value reflect.Value) reflect.Value {
	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}

		value = value.Elem()
	}

	return value
}

// SetValue set value to addressable point
//
//nolint:funlen,cyclop // normal length for this function and cyclomatic length
func SetValue(rvalue reflect.Value, value string) error {
	rvalue = Indirect(rvalue)
	rtype := rvalue.Type()

	if !rvalue.CanAddr() {
		return ErrCanAddr
	}

	addr := rvalue.Addr().Interface()

	if impl, ok := addr.(Setter); ok {
		return impl.Set(value) //nolint:wrapcheck // return error
	}

	if impl, ok := addr.(encoding.TextUnmarshaler); ok {
		return impl.UnmarshalText([]byte(value)) //nolint:wrapcheck // return error
	}

	if impl, ok := addr.(encoding.BinaryUnmarshaler); ok {
		return impl.UnmarshalBinary([]byte(value)) //nolint:wrapcheck // return error
	}

	//nolint:exhaustive // only needed types
	switch rtype.Kind() {
	case reflect.String:
		rvalue.SetString(value)

		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(value, 0, rtype.Bits())
		if err != nil {
			return err //nolint:wrapcheck // return error
		}

		rvalue.SetInt(val)

		break

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(value, 0, rtype.Bits())
		if err != nil {
			return err //nolint:wrapcheck // return error
		}

		rvalue.SetUint(val)

		break

	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(value, rtype.Bits())
		if err != nil {
			return err //nolint:wrapcheck // return error
		}

		rvalue.SetFloat(val)

		break

	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return err //nolint:wrapcheck // return error
		}

		rvalue.SetBool(val)

		break

	case reflect.Slice:
		if rtype.Elem().Kind() == reflect.Uint8 {
			rvalue.Set(reflect.ValueOf([]byte(value)))

			return nil
		}

		fallthrough
	default:
		return json.Unmarshal([]byte(value), addr) //nolint:wrapcheck // return error
	}

	return nil
}
