package http

import (
	"context"
	"net/http"
	"net/url"
	"reflect"

	"github.com/kamilov/go-kit/coder"
	"github.com/kamilov/go-kit/utils/structure"
)

type (
	requestDecoder[T any] func(ctx context.Context, data *T, request *http.Request) error
	pathData              http.Request
	valuesData            url.Values
	headerData            http.Header
)

const parseFormMaxMemory = 32 << 20

func (d *pathData) Get(key string) string {
	return (*http.Request)(d).PathValue(key)
}

func (d valuesData) Get(key string) string {
	return (url.Values)(d).Get(key)
}

func (d headerData) Get(key string) string {
	return (http.Header)(d).Get(key)
}

func newRequestDecoder[T any]() requestDecoder[T] {
	path, _ := structure.NewTypedDecoder[T]("path")
	query, _ := structure.NewTypedDecoder[T]("query")
	header, _ := structure.NewTypedDecoder[T]("header")
	form, _ := structure.NewTypedDecoder[T]("form")

	return decodeRequest[T](path, query, header, form)
}

//nolint:gocognit
func decodeRequest[T any](path, query, header, form *structure.TypedDecoder[T]) requestDecoder[T] {
	bodyDecoder := decodeBody[T]()

	return func(ctx context.Context, data *T, request *http.Request) error {
		if err := bodyDecoder(ctx, data, request); err != nil {
			return err
		}

		rv := reflect.Indirect(reflect.ValueOf(data))

		if path != nil {
			if err := path.DecodeValue((*pathData)(request), rv); err != nil {
				return err
			}
		}

		if query != nil {
			if err := query.DecodeValue((valuesData)(request.URL.Query()), rv); err != nil {
				return err
			}
		}

		if header != nil {
			if err := header.DecodeValue((headerData)(request.Header), rv); err != nil {
				return err
			}
		}

		if form != nil && !isEmptyBody(request) {
			if err := request.ParseMultipartForm(parseFormMaxMemory); err != nil {
				return err
			}

			if err := form.DecodeValue((valuesData)(request.Form), rv); err != nil {
				return err
			}
		}

		return nil
	}
}

func decodeBody[T any]() requestDecoder[T] {
	return func(ctx context.Context, data *T, request *http.Request) error {
		if isEmptyBody(request) {
			return nil
		}
		contentType := request.Header.Get("content-type")

		dec := coder.GetDecoder(getContentType([]byte(contentType)))
		if dec == nil {
			return ErrUnsupportedMediaType
		}

		return dec(ctx, request.Body, data)
	}
}

func isEmptyBody(request *http.Request) bool {
	return request.Header.Get("Content-Length") == "0" ||
		request.Method == http.MethodGet ||
		request.Method == http.MethodHead
}
