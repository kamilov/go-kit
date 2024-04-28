package structure

import (
	"encoding"
	"reflect"
	"strconv"
	"unsafe"
)

type (
	decoder        = func(reflect.Value, Data) error
	setFunc[T any] func(reflect.Value, T)
	getFunc        func(reflect.Value) reflect.Value
)

const (
	numDefaultSize = strconv.IntSize
	num8BitSize    = 8
	num16BitSize   = 16
	num32BitSize   = 32
	num64BitSize   = 64
)

var (
	textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
)

func noopDecoder(_ reflect.Value, _ Data) error { return nil }

func get(rt reflect.Type, index int, isPointer bool) getFunc {
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

func set[T any](rt reflect.Type, index int, isPointer bool) setFunc[T] {
	if isPointer {
		return func(rv reflect.Value, value T) {
			field := rv.Field(index)
			if field.IsNil() {
				field.Set(reflect.New(rt))
			}

			*(*T)(unsafe.Pointer(field.Elem().UnsafeAddr())) = value
		}
	}

	return func(rv reflect.Value, value T) {
		*(*T)(unsafe.Pointer(rv.Field(index).UnsafeAddr())) = value
	}
}

func decodeTextUnmarshaler(get getFunc, key string) decoder {
	return func(rv reflect.Value, data Data) error {
		if value := data.Get(key); value != "" {
			return get(rv).Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(value))
		}

		return nil
	}
}

func decodeString(set setFunc[string], key string) decoder {
	return func(rv reflect.Value, data Data) error {
		if value := data.Get(key); value != "" {
			set(rv, value)
		}

		return nil
	}
}

func decodeInt[T int | int8 | int16 | int32 | int64](set setFunc[T], key string, bits int) decoder {
	return func(rv reflect.Value, data Data) error {
		if value := data.Get(key); value != "" {
			i, err := strconv.ParseInt(value, 10, bits)
			if err != nil {
				return err
			}

			set(rv, T(i))
		}

		return nil
	}
}

func decodeUint[T uint | uint8 | uint16 | uint32 | uint64](set setFunc[T], key string, bits int) decoder {
	return func(rv reflect.Value, data Data) error {
		if value := data.Get(key); value != "" {
			i, err := strconv.ParseUint(value, 10, bits)
			if err != nil {
				return err
			}

			set(rv, T(i))
		}

		return nil
	}
}

func decodeFloat[T float32 | float64](set setFunc[T], key string, bits int) decoder {
	return func(rv reflect.Value, data Data) error {
		if value := data.Get(key); value != "" {
			f, err := strconv.ParseFloat(value, bits)
			if err != nil {
				return err
			}

			set(rv, T(f))
		}

		return nil
	}
}

func decodeBool(set setFunc[bool], key string) decoder {
	return func(rv reflect.Value, data Data) error {
		if value := data.Get(key); value != "" {
			b, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}

			set(rv, b)
		}

		return nil
	}
}

func decodeBytes(set setFunc[[]byte], key string) decoder {
	return func(rv reflect.Value, data Data) error {
		if value := data.Get(key); value != "" {
			set(rv, []byte(value))
		}

		return nil
	}
}

func compileDecoders(srt reflect.Type, tagName string, isPointer bool) (decoder, error) {
	var decoders []decoder

	for i := 0; i < srt.NumField(); i++ {
		field := srt.Field(i)
		if field.PkgPath != "" {
			continue
		}

		fieldType, fieldKind, fieldIsPointer := typeKind(field.Type)

		tagValue, ok := field.Tag.Lookup(tagName)
		if !ok && fieldKind != reflect.Struct {
			continue
		}

		if reflect.PointerTo(fieldType).Implements(textUnmarshalerType) {
			decoders = append(decoders, decodeTextUnmarshaler(get(fieldType, i, fieldIsPointer), tagValue))
			continue
		}

		switch fieldKind {
		case reflect.Struct:
			dec, err := compileDecoders(fieldType, tagName, fieldIsPointer)
			if err != nil {
				return nil, err
			}
			index := i
			decoders = append(decoders, func(rv reflect.Value, data Data) error {
				return dec(rv.Field(index), data)
			})

		case reflect.String:
			decoders = append(decoders, decodeString(set[string](fieldType, i, fieldIsPointer), tagValue))

		case reflect.Bool:
			decoders = append(decoders, decodeBool(set[bool](fieldType, i, fieldIsPointer), tagValue))

		case reflect.Int:
			decoders = append(decoders, decodeInt(set[int](fieldType, i, fieldIsPointer), tagValue, numDefaultSize))

		case reflect.Int8:
			decoders = append(decoders, decodeInt(set[int8](fieldType, i, isPointer), tagValue, num8BitSize))

		case reflect.Int16:
			decoders = append(decoders, decodeInt(set[int16](fieldType, i, isPointer), tagValue, num16BitSize))

		case reflect.Int32:
			decoders = append(decoders, decodeInt(set[int32](fieldType, i, isPointer), tagValue, num32BitSize))

		case reflect.Int64:
			decoders = append(decoders, decodeInt(set[int64](fieldType, i, isPointer), tagValue, num64BitSize))

		case reflect.Uint:
			decoders = append(decoders, decodeUint(set[uint](fieldType, i, isPointer), tagValue, numDefaultSize))

		case reflect.Uint8:
			decoders = append(decoders, decodeUint(set[uint8](fieldType, i, isPointer), tagValue, num8BitSize))

		case reflect.Uint16:
			decoders = append(decoders, decodeUint(set[uint16](fieldType, i, isPointer), tagValue, num16BitSize))

		case reflect.Uint32:
			decoders = append(decoders, decodeUint(set[uint32](fieldType, i, isPointer), tagValue, num32BitSize))

		case reflect.Uint64:
			decoders = append(decoders, decodeUint(set[uint64](fieldType, i, isPointer), tagValue, num64BitSize))

		case reflect.Float32:
			decoders = append(decoders, decodeFloat(set[float32](fieldType, i, isPointer), tagValue, num32BitSize))

		case reflect.Float64:
			decoders = append(decoders, decodeFloat(set[float64](fieldType, i, isPointer), tagValue, num64BitSize))

		case reflect.Slice:
			if fieldType.Elem().Kind() == reflect.Uint8 {
				decoders = append(decoders, decodeBytes(set[[]byte](fieldType, i, isPointer), tagValue))
				continue
			}
			fallthrough

		default:
			return nil, ErrUnsupportedType
		}
	}

	if len(decoders) == 0 {
		return nil, ErrTagNotFound
	}

	return func(rv reflect.Value, data Data) error {
		if isPointer {
			if rv.IsNil() {
				rv.Set(reflect.New(srt))
			}

			rv = rv.Elem()
		}

		for _, dec := range decoders {
			if err := dec(rv, data); err != nil {
				return err
			}
		}

		return nil
	}, nil
}
