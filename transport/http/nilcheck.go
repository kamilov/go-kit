package http

import (
	"reflect"
)

func NewNilCheck(zero any) func(value any) bool {
	rt := reflect.TypeOf(zero)

	if rt == nil {
		return func(value any) bool { return value == nil }
	}

	//nolint:exhaustive // use only currently types
	switch rt.Kind() {
	case reflect.String, reflect.Pointer, reflect.Interface:
		return func(value any) bool { return value == zero }

	case reflect.Map:
		return func(value any) bool { return dataOf(value) == nil }

	case reflect.Slice:
		return func(value any) bool {
			//nolint:staticcheck // we are using this approach for now
			return (*reflect.SliceHeader)(dataOf(value)).Data == 0
		}

	default:
		return func(_ any) bool { return false }
	}
}
