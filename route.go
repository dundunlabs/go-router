package gorouter

import "net/http"

type Handler func(http.ResponseWriter, *http.Request)

type Route struct {
	Path     string
	Method   string
	Handler  Handler
	Children []Route
}
