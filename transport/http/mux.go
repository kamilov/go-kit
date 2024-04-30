package http

import (
	"net/http"
	"strings"

	"github.com/kamilov/go-kit/endpoint"
)

func Group(server *Server, path string) *Server {
	s := *server
	newServer := &s
	newServer.basePath += strings.Trim(path, "/") + "/"

	return newServer
}

func Use(server *Server, middlewares ...Middleware) {
	server.middlewares = append(server.middlewares, middlewares...)
}

func Get[Input, Output any](
	server *Server,
	path string,
	controller endpoint.Endpoint[Input, Output],
	middlewares ...endpoint.Middleware[Input, Output],
) {
	register(server, http.MethodGet, path, controller, middlewares...)
}

func Post[Input, Output any](
	server *Server,
	path string,
	controller endpoint.Endpoint[Input, Output],
	middlewares ...endpoint.Middleware[Input, Output],
) {
	register(server, http.MethodPost, path, controller, middlewares...)
}

func Put[Input, Output any](
	server *Server,
	path string,
	controller endpoint.Endpoint[Input, Output],
	middlewares ...endpoint.Middleware[Input, Output],
) {
	register(server, http.MethodPut, path, controller, middlewares...)
}

func Patch[Input, Output any](
	server *Server,
	path string,
	controller endpoint.Endpoint[Input, Output],
	middlewares ...endpoint.Middleware[Input, Output],
) {
	register(server, http.MethodPatch, path, controller, middlewares...)
}

func Delete[Input, Output any](
	server *Server,
	path string,
	controller endpoint.Endpoint[Input, Output],
	middlewares ...endpoint.Middleware[Input, Output],
) {
	register(server, http.MethodDelete, path, controller, middlewares...)
}

func register[Input, Output any](
	server *Server,
	method, path string,
	controller endpoint.Endpoint[Input, Output],
	middlewares ...endpoint.Middleware[Input, Output],
) {
	pattern := method + " " + server.basePath + strings.TrimLeft(path, "/")
	controller = withEndpointMiddlewares(controller, middlewares...)
	muxHandler := handler(server, controller)
	muxHandler = withHTTPMiddlewares(muxHandler, server.middlewares...)

	server.mux.Handle(pattern, muxHandler)
}
