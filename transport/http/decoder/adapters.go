package decoder

import (
	"net/http"
	"net/url"
)

type (
	Path struct {
		request *http.Request
	}
	Query   url.Values
	Headers http.Header
)

func (p Path) Get(key string) string {
	return p.request.PathValue(key)
}

func (q Query) Get(key string) string {
	return q.Get(key)
}

func (h Headers) Get(key string) string {
	return h.Get(key)
}
