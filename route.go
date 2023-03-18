package gorouter

import "net/http"

type HandlerFunc func(http.ResponseWriter, Request)
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

type Route struct {
	Path       string
	Method     string
	Handler    HandlerFunc
	Middleware MiddlewareFunc
	Children   []Route
}
