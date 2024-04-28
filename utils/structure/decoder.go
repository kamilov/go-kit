package structure

import (
	"errors"
	"reflect"
	"sync"
)

type (
	Data interface {
		Get(string) string
	}

	Decoder struct {
		tag   string
		cache sync.Map
	}

	TypedDecoder[T any] struct {
		dec decoder
	}
)

var (
	ErrUnsupportedType = errors.New("unsupported type")
	ErrTagNotFound     = errors.New("tag not found")
)

func NewDecoder(tag string) *Decoder {
	return &Decoder{tag: tag}
}

func NewTypedDecoder[T any](tag string) (*TypedDecoder[T], error) {
	var target T

	rt := reflect.TypeOf(target)
	if rt == nil {
		return nil, ErrUnsupportedType
	}

	rt, rtk, isPointer := typeKind(rt)

	if rtk != reflect.Struct {
		return nil, ErrStructPointer
	}

	dec, err := compileDecoders(rt, tag, isPointer)
	if err != nil {
		return nil, err
	}

	return &TypedDecoder[T]{dec}, nil
}

func (d *Decoder) Decode(data Data, target any) error {
	if err := ValidateStructPointer(target); err != nil {
		return err
	}

	rv := reflect.Indirect(reflect.ValueOf(target))
	rt := rv.Type()

	cacheKey := rv.Type()
	dec, ok := d.cache.Load(cacheKey)

	if !ok {
		var err error

		dec, err = compileDecoders(rt, d.tag, rt.Kind() == reflect.Pointer)
		if err != nil {
			if !errors.Is(err, ErrTagNotFound) {
				return err
			}

			dec = noopDecoder
		}

		d.cache.Store(cacheKey, dec)
	}

	return dec.(decoder)(rv, data)
}

func (d *TypedDecoder[T]) Decode(data Data, target *T) error {
	return d.DecodeValue(data, reflect.ValueOf(target).Elem())
}

func (d *TypedDecoder[T]) DecodeValue(data Data, rv reflect.Value) error {
	return d.dec(rv, data)
}
