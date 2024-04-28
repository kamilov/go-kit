package decoder

import (
	"encoding"
	"errors"
	"reflect"
	"strconv"
	"unsafe"
)

type (
	Adapter interface {
		Get(string) string
	}

	decoder = func(reflect.Value, Adapter) error
)

var (
	ErrUnsupportedType = errors.New("unsupported type")
	ErrTagNotFound     = errors.New("tag not found")

	textUnmarshalerType   = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
	binaryUnmarshalerType = reflect.TypeOf((*encoding.BinaryUnmarshaler)(nil)).Elem()
)

func compile(rt reflect.Type, tagName string) (decoder, error) {
	var decoders []decoder

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.PkgPath != "" {
			continue
		}

		fieldType, fieldKind, fieldIsPointer := typeKind(field.Type)

		tag, ok := field.Tag.Lookup(tagName)
		if !ok && fieldKind != reflect.Struct {
			continue
		}

		if reflect.PointerTo(fieldType).Implements(textUnmarshalerType) {
			decoders = append(decoders, decoderTextUnmarshaler(get(fieldType, i, fieldIsPointer), tag))
			continue
		}

		if reflect.PointerTo(fieldType).Implements(binaryUnmarshalerType) {
			decoders = append(decoders, decoderBinaryUnmarshaler(get(fieldType, i, fieldIsPointer), tag))
			continue
		}

		switch fieldKind {
		case reflect.Struct:
			dec, err := compile(field.Type, tagName)
			if err != nil {
				return nil, err
			}
			index := i
			decoders = append(decoders, func(rv reflect.Value, adapter Adapter) error {
				return dec(rv.Field(index), adapter)
			})

		case reflect.String:
			decoders = append(decoders, decoderString(set[string](fieldType, i, fieldIsPointer), tag))

		case reflect.Int:
			decoders = append(decoders, decoderInt(set[int](fieldType, i, fieldIsPointer), tag, strconv.IntSize))

		case reflect.Int8:
			decoders = append(decoders, decoderInt(set[int8](fieldType, i, fieldIsPointer), tag, 8))

		case reflect.Int16:
			decoders = append(decoders, decoderInt(set[int16](fieldType, i, fieldIsPointer), tag, 16))

		case reflect.Int32:
			decoders = append(decoders, decoderInt(set[int32](fieldType, i, fieldIsPointer), tag, 32))

		case reflect.Int64:
			decoders = append(decoders, decoderInt(set[int64](fieldType, i, fieldIsPointer), tag, 64))

		case reflect.Uint:
			decoders = append(decoders, decoderUint(set[uint](fieldType, i, fieldIsPointer), tag, strconv.IntSize))

		case reflect.Uint8:
			decoders = append(decoders, decoderUint(set[uint8](fieldType, i, fieldIsPointer), tag, 8))

		case reflect.Uint16:
			decoders = append(decoders, decoderUint(set[uint16](fieldType, i, fieldIsPointer), tag, 16))

		case reflect.Uint32:
			decoders = append(decoders, decoderUint(set[uint32](fieldType, i, fieldIsPointer), tag, 32))

		case reflect.Uint64:
			decoders = append(decoders, decoderUint(set[uint64](fieldType, i, fieldIsPointer), tag, 64))

		case reflect.Float32:
			decoders = append(decoders, decodeFloat(set[float32](fieldType, i, fieldIsPointer), tag, 32))

		case reflect.Float64:
			decoders = append(decoders, decodeFloat(set[float64](fieldType, i, fieldIsPointer), tag, 64))

		case reflect.Bool:
			decoders = append(decoders, decoderBool(set[bool](fieldType, i, fieldIsPointer), tag))

		case reflect.Slice:
			decoders = append(decoders, decoderBytes(set[[]byte](fieldType, i, fieldIsPointer), tag))

		default:
			return nil, ErrUnsupportedType
		}
	}

	if len(decoders) == 0 {
		return nil, ErrTagNotFound
	}

	return func(rv reflect.Value, adapter Adapter) error {
		if rt.Kind() == reflect.Pointer {
			if rv.IsNil() {
				rv.Set(reflect.New(rt))
			}

			rv = rv.Elem()
		}

		for _, dec := range decoders {
			if err := dec(rv, adapter); err != nil {
				return err
			}
		}

		return nil
	}, nil
}

func typeKind(rt reflect.Type) (reflect.Type, reflect.Kind, bool) {
	if rt.Kind() == reflect.Pointer {
		return rt.Elem(), rt.Elem().Kind(), true
	}
	return rt, rt.Kind(), false
}

func set[T any](rt reflect.Type, index int, isPointer bool) func(reflect.Value, T) {
	if isPointer {
		return func(rv reflect.Value, t T) {
			field := rv.Field(index)
			if field.IsNil() {
				field.Set(reflect.New(rt))
			}

			*(*T)(unsafe.Pointer(field.Elem().UnsafeAddr())) = t
		}
	}

	return func(rv reflect.Value, t T) {
		*(*T)(unsafe.Pointer(rv.Field(index).UnsafeAddr())) = t
	}
}

func get(rt reflect.Type, index int, isPointer bool) func(reflect.Value) reflect.Value {
	if isPointer {
		return func(rv reflect.Value) reflect.Value {
			field := rv.Field(index)
			if field.IsNil() {
				field.Set(reflect.New(rt))
			}
			return field
		}
	}

	return func(rv reflect.Value) reflect.Value {
		return rv.Field(index).Addr()
	}
}

func decoderTextUnmarshaler(get func(reflect.Value) reflect.Value, key string) decoder {
	return func(rv reflect.Value, adapter Adapter) error {
		if value := adapter.Get(key); value != "" {
			return get(rv).Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(value))
		}

		return nil
	}
}

func decoderBinaryUnmarshaler(get func(reflect.Value) reflect.Value, key string) decoder {
	return func(rv reflect.Value, adapter Adapter) error {
		if value := adapter.Get(key); value != "" {
			return get(rv).Interface().(encoding.BinaryUnmarshaler).UnmarshalBinary([]byte(value))
		}

		return nil
	}
}

func decoderString(set func(reflect.Value, string), key string) decoder {
	return func(rv reflect.Value, adapter Adapter) error {
		if value := adapter.Get(key); value != "" {
			set(rv, value)
		}

		return nil
	}
}

func decoderInt[T int | int8 | int16 | int32 | int64](set func(reflect.Value, T), key string, bits int) decoder {
	return func(rv reflect.Value, adapter Adapter) error {
		if value := adapter.Get(key); value != "" {
			i, err := strconv.ParseInt(value, 10, bits)
			if err != nil {
				return err
			}

			set(rv, T(i))
		}

		return nil
	}
}

func decoderUint[T uint | uint8 | uint16 | uint32 | uint64](set func(reflect.Value, T), key string, bits int) decoder {
	return func(rv reflect.Value, adapter Adapter) error {
		if value := adapter.Get(key); value != "" {
			i, err := strconv.ParseUint(value, 10, bits)
			if err != nil {
				return err
			}

			set(rv, T(i))
		}

		return nil
	}
}

func decodeFloat[T float32 | float64](set func(reflect.Value, T), key string, bits int) decoder {
	return func(rv reflect.Value, adapter Adapter) error {
		if value := adapter.Get(key); value != "" {
			f, err := strconv.ParseFloat(value, bits)
			if err != nil {
				return err
			}

			set(rv, T(f))
		}

		return nil
	}
}

func decoderBool(set func(reflect.Value, bool), key string) decoder {
	return func(rv reflect.Value, adapter Adapter) error {
		if value := adapter.Get(key); value != "" {
			b, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}

			set(rv, b)
		}

		return nil
	}
}

func decoderBytes(set func(reflect.Value, []byte), key string) decoder {
	return func(rv reflect.Value, adapter Adapter) error {
		if value := adapter.Get(key); value != "" {
			set(rv, []byte(value))
		}

		return nil
	}
}
