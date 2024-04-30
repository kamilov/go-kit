package structure_test

import (
	"errors"
	"reflect"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kamilov/go-kit/utils/structure"
)

type unmarshalText string
type Map map[string]string

func (m Map) Get(key string) string {
	return m[key]
}

func (u *unmarshalText) UnmarshalText(text []byte) error {
	*u = unmarshalText("_t_" + string(text) + "_t_")
	return nil
}

func TestDecoder(t *testing.T) {
	type child struct {
		String string `field:"string"`
	}

	type test struct {
		UnmarshalText        unmarshalText  `field:"string"`
		UnmarshalTextPointer *unmarshalText `field:"string"`
		unexported           string         `field:"string"` //nolint:unused // for only test
		String               string         `field:"string"`
		StringPointer        *string        `field:"string"`
		Int                  int            `field:"number"`
		Int8                 int8           `field:"number"`
		Int16                int16          `field:"number"`
		Int32                int32          `field:"number"`
		Int64                int64          `field:"number"`
		Uint                 uint           `field:"number"`
		Uint8                uint8          `field:"number"`
		Uint16               uint16         `field:"number"`
		Uint32               uint32         `field:"number"`
		Uint64               uint64         `field:"number"`
		Float32              float32        `field:"number"`
		Float64              float64        `field:"number"`
		Bool                 bool           `field:"bool"`
		Bytes                []byte         `field:"string"`
		Child                child
		ChildPointer         *child
	}

	testMap := Map{
		"string": "string",
		"number": "1",
		"bool":   "TRUE",
	}

	str := "string"
	text := unmarshalText("_t_string_t_")
	expected := &test{
		UnmarshalText:        text,
		UnmarshalTextPointer: &text,
		String:               str,
		StringPointer:        &str,
		Int:                  1,
		Int8:                 1,
		Int16:                1,
		Int32:                1,
		Int64:                1,
		Uint:                 1,
		Uint8:                1,
		Uint16:               1,
		Uint32:               1,
		Uint64:               1,
		Float32:              1,
		Float64:              1,
		Bool:                 true,
		Bytes:                []byte("string"),
		Child: child{
			String: "string",
		},
		ChildPointer: &child{
			String: "string",
		},
	}

	exportAll := cmp.Exporter(func(_ reflect.Type) bool { return true })

	t.Run("Decoder", func(t *testing.T) {
		t.Helper()

		dec := structure.NewDecoder("field")
		data := &test{}

		if err := dec.Decode(testMap, data); err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(expected, data, exportAll); diff != "" {
			t.Errorf(diff)
		}
	})

	t.Run("TypedDecoder", func(t *testing.T) {
		t.Helper()

		dec, err := structure.NewTypedDecoder[test]("field")
		if err != nil {
			t.Fatal(err)
		}

		data := &test{}

		if err = dec.Decode(testMap, data); err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(expected, data, exportAll); diff != "" {
			t.Errorf(diff)
		}

		data = &test{}
		rv := reflect.ValueOf(data).Elem()

		if err = dec.DecodeValue(testMap, rv); err != nil {
			t.Fatal(err)
		}

		if diff := cmp.Diff(expected, data, exportAll); diff != "" {
			t.Errorf(diff)
		}
	})
}

func TestDecoderErrors(t *testing.T) {
	testMap := Map{"test": "string"}

	tests := []struct {
		data any
		err  error
	}{
		{
			nil,
			structure.ErrStructPointer,
		},
		{
			"",
			structure.ErrStructPointer,
		},
		{
			new(string),
			structure.ErrStructPointer,
		},
		{
			new(int),
			structure.ErrStructPointer,
		},
		{
			&struct {
				Test string `json:"test"`
			}{},
			structure.ErrTagNotFound,
		},
		{
			&struct {
				Test chan string `field:"test"`
			}{},
			structure.ErrUnsupportedType,
		},
		{
			&struct {
				Test  string `field:"test"`
				Child struct {
					Test chan string `field:"test"`
				}
			}{},
			structure.ErrUnsupportedType,
		},
		{
			&struct {
				Test int `field:"test"`
			}{},
			strconv.ErrSyntax,
		},
		{
			&struct {
				Test float64 `field:"test"`
			}{},
			strconv.ErrSyntax,
		},
		{
			&struct {
				Test bool `field:"test"`
			}{},
			strconv.ErrSyntax,
		},
	}

	t.Run("Decoder", func(t *testing.T) {
		for _, test := range tests {
			dec := structure.NewDecoder("field")
			err := dec.Decode(testMap, test.data)

			if errors.Is(test.err, structure.ErrTagNotFound) {
				if err != nil {
					t.Errorf("should silently ignore error %q for %T", test.err, test.data)
				}
			} else {
				if !errors.Is(err, test.err) {
					t.Errorf("should return %q for %T: %q", test.err, test.data, err)
				}
			}
		}
	})
}
