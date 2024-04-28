package structure

import (
	"errors"
	"reflect"
)

var (
	ErrPointer       = errors.New("must be a pointer")
	ErrStructPointer = errors.New("must be a pointer to a struct")
	ErrNilPointer    = errors.New("the pointer should not be nil")
)

func ValidatePointer(pointer any) error {
	rv := reflect.ValueOf(pointer)

	if rv.Kind() != reflect.Ptr {
		return ErrPointer
	}

	return nil
}

func ValidateStructPointer(structPointer any) error {
	rv := reflect.ValueOf(structPointer)

	if rv.Kind() != reflect.Pointer || !rv.IsNil() && rv.Elem().Kind() != reflect.Struct {
		return ErrStructPointer
	}

	if rv.IsNil() {
		return ErrNilPointer
	}

	return nil
}
