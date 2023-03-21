package gorouter

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type Request struct {
	*http.Request
	node   *node
	params Params
	body   Body
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

func (r *Request) Param(key string) string {
	return r.Params()[key]
}

func (r *Request) ParseBody() (Body, error) {
	var err error
	if r.body == nil {
		r.body, err = io.ReadAll(r.Body)
	}
	return r.body, err
}

func (r *Request) MustParseBody() Body {
	v, err := r.ParseBody()
	if err != nil {
		panic(err)
	}
	return v
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

type Body []byte

func (r Body) Bind(v any) error {
	return json.Unmarshal(r, v)
}

func (r Body) MustBind(v any) {
	if err := json.Unmarshal(r, v); err != nil {
		panic(err)
	}
}
