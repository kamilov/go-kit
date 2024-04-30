package http

import (
	"bytes"
	"unsafe"
)

type emptyInterface struct {
	typ unsafe.Pointer
	ptr unsafe.Pointer
}

func dataOf(value any) unsafe.Pointer {
	return (*emptyInterface)(unsafe.Pointer(&value)).ptr
}

func getContentType(contentType []byte) string {
	index := bytes.IndexRune(contentType, ';')
	if index > 0 {
		contentType = contentType[:index]
	}

	return string(bytes.TrimSpace(contentType))
}
