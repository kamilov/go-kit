package text

import (
	"context"
	"encoding"
	"fmt"
	"io"
	"strconv"
)

func encoder(_ context.Context, writer io.Writer, data any) error {
	if data == nil {
		return nil
	}

	var (
		bytes []byte
		err   error
	)

	switch data := data.(type) {
	case string:
		bytes = []byte(data)

	case []byte:
		bytes = data

	case int:
		bytes = strconv.AppendInt(bytes, int64(data), numBase)

	case int8:
		bytes = strconv.AppendInt(bytes, int64(data), numBase)

	case int16:
		bytes = strconv.AppendInt(bytes, int64(data), numBase)

	case int32:
		bytes = strconv.AppendInt(bytes, int64(data), numBase)

	case int64:
		bytes = strconv.AppendInt(bytes, data, numBase)

	case uint:
		bytes = strconv.AppendUint(bytes, uint64(data), numBase)

	case uint8:
		bytes = strconv.AppendUint(bytes, uint64(data), numBase)

	case uint16:
		bytes = strconv.AppendUint(bytes, uint64(data), numBase)

	case uint32:
		bytes = strconv.AppendUint(bytes, uint64(data), numBase)

	case uint64:
		bytes = strconv.AppendUint(bytes, data, numBase)

	case float32:
		bytes = strconv.AppendFloat(bytes, float64(data), 'f', -1, num32BitSize)

	case float64:
		bytes = strconv.AppendFloat(bytes, data, 'f', -1, num64BitSize)

	case bool:
		bytes = strconv.AppendBool(bytes, data)

	case encoding.TextMarshaler:
		bytes, err = data.MarshalText()

	case error:
		bytes = []byte(data.Error())

	case fmt.Stringer:
		bytes = []byte(data.String())

	default:
		return ErrUnsupportedType
	}

	if err != nil {
		return err
	}

	_, _ = writer.Write(bytes)

	return nil
}
