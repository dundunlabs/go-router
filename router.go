package gorouter

import (
	"net/http"
)

func New(routes []Route) Router {
	tree := newNode(nil)
	tree.generateFromRoutes(routes, "", nil)
	return Router{tree}
}

type Router struct {
	tree *node
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.EscapedPath()
	n := router.tree.findNode(path)

	if n == nil || len(n.handlers) == 0 {
		http.NotFound(w, r)
		return
	}

	handle, ok := n.handlers[r.Method]

	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	req := Request{
		Request: r,
		node:    n,
	}
	handle(w, req)
}
