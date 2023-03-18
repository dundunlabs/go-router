package gorouter

import (
	"context"
	"net/http"
)

func New(routes []Route) Router {
	tree := newNode("", nil)
	tree.generateFromRoutes(routes, nil)
	return Router{tree}
}

type Router struct {
	tree *node
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.EscapedPath()
	n := router.tree.findNode(path)

	if n == nil {
		http.NotFound(w, r)
		return
	}

	handle, ok := n.handlers[r.Method]

	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	handle(w, r.WithContext(context.WithValue(r.Context(), "route", n.route)))
}
