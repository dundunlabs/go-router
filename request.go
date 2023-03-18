package gorouter

import (
	"net/http"
	"strings"
)

type Request struct {
	*http.Request
	node   *node
	params Params
}

type Params map[string]string

func (p Params) parse(route string, path string) {
	if i := strings.IndexByte(route, ':'); i > -1 {
		if j1 := strings.IndexByte(route[i+1:], '/') + i + 1; j1 > i+1 {
			if j2 := strings.IndexByte(path[i:], '/') + i; j2 > i {
				p[route[i+1:j1]] = path[i:j2]
				p.parse(route[j1+1:], path[j2+1:])
			}
		} else {
			p[route[i+1:]] = path[i:]
		}
	}
}

func (r *Request) Route() string {
	return r.node.route
}

func (r *Request) Params() Params {
	route := r.Route()

	if strings.IndexByte(route, ':') == -1 {
		return nil
	}

	if r.params != nil {
		return r.params
	}

	r.params = make(Params)
	r.params.parse(route, r.URL.EscapedPath())

	return r.params
}
