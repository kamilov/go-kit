package structure

import "reflect"

func typeKind(rt reflect.Type) (reflect.Type, reflect.Kind, bool) {
	if rt.Kind() == reflect.Pointer {
		return rt.Elem(), rt.Elem().Kind(), true
	}

	return rt, rt.Kind(), false
}
