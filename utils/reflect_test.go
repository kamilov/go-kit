package utils

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	testSet    bool
	testInt    int64
	testString string
)

func (t *testSet) Set(value string) error {
	val, err := strconv.ParseBool(value)

	if err != nil {
		return err
	}

	*t = testSet(val)

	return nil
}

func (t *testInt) UnmarshalText(data []byte) error {
	val, err := strconv.ParseUint(string(data), 10, 0)

	if err != nil {
		return err
	}

	*t = testInt(val)

	return nil
}

func (t *testString) UnmarshalBinary(data []byte) error {
	*t = testString(string(data) + "-ok")
	return nil
}

func TestIndirect(t *testing.T) {
	var a int

	assert.Equal(t, reflect.ValueOf(a).Kind(), Indirect(reflect.ValueOf(a)).Kind())

	var b *int

	assert.Equal(t, reflect.ValueOf(a).Kind(), Indirect(reflect.ValueOf(&b)).Kind())
	assert.NotNil(t, b)
	assert.Equal(t, 0, *b)
}

func TestSetValue(t *testing.T) {
	cfg := struct {
		str           string
		strPtr        *string
		int1          int
		uint1         uint
		float1        float32
		bool1         bool
		slice1        []byte
		slice2        []int
		slice3        []string
		map1          map[string]int
		testInt       testInt
		testIntPtr    *testInt
		testString    testString
		testStringPtr *testString
		testSet       testSet
	}{}

	tests := []struct {
		name     string
		rvalue   reflect.Value
		value    string
		expected any
		isEqual  bool
		isError  bool
	}{
		{"Test Str", reflect.ValueOf(&cfg.str), "test", "test", true, false},
		{"Test Str Pointer", reflect.ValueOf(&cfg.strPtr), "test", "test", true, false},
		{"Test Int", reflect.ValueOf(&cfg.int1), "1", int(1), true, false},
		{"Test Int Not equal", reflect.ValueOf(&cfg.int1), "1", uint(1), false, false},
		{"test Int Error", reflect.ValueOf(&cfg.int1), "a1", int(1), true, true},
		{"Test Uint", reflect.ValueOf(&cfg.uint1), "1", uint(1), true, false},
		{"Test Uint Not equal", reflect.ValueOf(&cfg.uint1), "1", int(1), false, false},
		{"Test Uint Error", reflect.ValueOf(&cfg.uint1), "a1", uint(1), true, true},
		{"Test Float", reflect.ValueOf(&cfg.float1), "10.1", float32(10.1), true, false},
		{"Test Float Not equal", reflect.ValueOf(&cfg.float1), "10.1", float64(10.1), false, false},
		{"Test Float Error", reflect.ValueOf(&cfg.float1), "a10.1", float32(10.1), true, true},
		{"Test Bool Int", reflect.ValueOf(&cfg.bool1), "1", true, true, false},
		{"Test Bool Text", reflect.ValueOf(&cfg.bool1), "TRUE", true, true, false},
		{"Test Bool Error", reflect.ValueOf(&cfg.bool1), "tRuE", true, true, true},
		{"Test Slice Bytes", reflect.ValueOf(&cfg.slice1), "abc", []byte("abc"), true, false},
		{"Test Slice Int", reflect.ValueOf(&cfg.slice2), "[1, 2]", []int{1, 2}, true, false},
		{"Test Slice DSN", reflect.ValueOf(&cfg.slice3), "[\"a\", \"b\", \"c\"]", []string{"a", "b", "c"}, true, false},
		{"Test Map", reflect.ValueOf(&cfg.map1), "{\"a\": 1, \"b\": 2, \"c\": 3}", map[string]int{"a": 1, "b": 2, "c": 3}, true, false},
		{"Test Map Error", reflect.ValueOf(&cfg.map1), "\"a\": 1, \"b\": 2, \"c\": 3", map[string]int{"a": 1, "b": 2, "c": 3}, true, true},
		{"Test TestInt", reflect.ValueOf(&cfg.testInt), "1", testInt(1), true, false},
		{"Test TestInt", reflect.ValueOf(&cfg.testIntPtr), "1", testInt(1), true, false},
		{"Test TestString", reflect.ValueOf(&cfg.testString), "test", testString("test-ok"), true, false},
		{"Test TestString", reflect.ValueOf(&cfg.testStringPtr), "test", testString("test-ok"), true, false},
		{"Test TestSet", reflect.ValueOf(&cfg.testSet), "1", testSet(true), true, false},
		{"Test TestSet", reflect.ValueOf(&cfg.testSet), "TRUE", testSet(true), true, false},
		{"Test TestSet", reflect.ValueOf(&cfg.testSet), "tRuE", testSet(true), true, true},
		{"Test Pointer Error", reflect.ValueOf("test"), "test", "test", true, true},
	}

	for _, test := range tests {
		err := SetValue(test.rvalue, test.value)

		if test.isError {
			assert.NotNil(t, err, test.name)
		} else if assert.Nil(t, err, test.name) {
			actual := Indirect(test.rvalue)

			if test.isEqual {
				assert.Equal(t, test.expected, actual.Interface(), test.name)
				assert.True(t, reflect.DeepEqual(test.expected, actual.Interface()), test.name)
			} else {
				assert.False(t, reflect.DeepEqual(test.expected, actual.Interface()), test.name)
			}
		}
	}
}
