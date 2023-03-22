package gorouter

import (
	"fmt"
	"net/http"
	"runtime/debug"
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

	if n == nil {
		http.NotFound(w, r)
		return
	}

	handle, ok := n.handlers[r.Method]

	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		}
	}()

	req := &Request{
		Request: r,
		node:    n,
	}
	res := &Response{
		ResponseWriter: w,
	}

	handle(req, res)

	if res.statusCode > 0 {
		w.WriteHeader(res.statusCode)
	}
	w.Write(res.b)
}
