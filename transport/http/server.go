package http

import (
	"context"
	"net/http"
	"time"

	"github.com/kamilov/go-kit/transport/http/content"
)

type Server struct {
	server *http.Server
	mux    *http.ServeMux

	basePath    string
	middlewares []Middleware

	negotiateTypes []content.ContentType
}

const (
	defaultAddr                   = ":9999"
	defaultBasePath               = "/"
	defaultTimeout  time.Duration = 30 * time.Second
)

func New(opts ...Option) *Server {
	s := &Server{
		server: &http.Server{
			Addr:              defaultAddr,
			ReadTimeout:       defaultTimeout,
			ReadHeaderTimeout: defaultTimeout,
			WriteTimeout:      defaultTimeout,
			IdleTimeout:       defaultTimeout,
		},
		mux: http.NewServeMux(),

		basePath:    defaultBasePath,
		middlewares: make([]Middleware, 0),
	}

	for _, opt := range opts {
		opt.apply(s)
	}

	return s
}

func (s *Server) Negotiate(types ...content.ContentType) *Server {
	s.negotiateTypes = types

	return s
}

func (s *Server) Run() error {
	s.server.Handler = s.mux

	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
