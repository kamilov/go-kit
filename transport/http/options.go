package http

import (
	"fmt"
	"strings"
	"time"

	"github.com/kamilov/go-kit/transport/http/content"
)

type (
	optionFunc func(server *Server)

	Option interface {
		apply(server *Server)
	}
)

func (f optionFunc) apply(server *Server) {
	f(server)
}

func WithPort(port int) Option {
	return optionFunc(func(server *Server) {
		server.server.Addr = fmt.Sprintf(":%d", port)
	})
}

func WithBasePath(path string) Option {
	return optionFunc(func(server *Server) {
		server.basePath = strings.TrimRight(path, "/") + "/"
	})
}

func WithNegotiateTypes(types ...content.ContentType) Option {
	return optionFunc(func(server *Server) {
		server.negotiateTypes = types
	})
}

func WithReadTimeout(timeout int64) Option {
	return optionFunc(func(server *Server) {
		server.server.ReadTimeout = time.Duration(timeout) * time.Second
	})
}

func WithReadHeaderTimeout(timeout int64) Option {
	return optionFunc(func(server *Server) {
		server.server.ReadHeaderTimeout = time.Duration(timeout) * time.Second
	})
}

func WithWriteTimeout(timeout int64) Option {
	return optionFunc(func(server *Server) {
		server.server.WriteTimeout = time.Duration(timeout) * time.Second
	})
}

func WithIdleTimeout(timeout int64) Option {
	return optionFunc(func(server *Server) {
		server.server.IdleTimeout = time.Duration(timeout) * time.Second
	})
}
