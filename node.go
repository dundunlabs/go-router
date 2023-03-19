package gorouter

import (
	"strings"
)

const paramKey = "?"

func newNode(parent *node) *node {
	return &node{
		parent:   parent,
		children: make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

type node struct {
	route    string
	parent   *node
	wildcard bool
	children map[string]*node
	handlers map[string]HandlerFunc
}

func (n *node) generateFromRoutes(routes []Route, prefix string, middleware MiddlewareFunc) {
	for _, r := range routes {
		cn := n.generateFromPath(r.Path)
		route := ""
		if prefix != "" && r.Path != "" {
			route = prefix + r.Path
		} else if prefix != "" {
			route = prefix
		} else if r.Path != "" {
			route = r.Path
		}

		var m MiddlewareFunc
		if middleware != nil && r.Middleware != nil {
			m = func(next HandlerFunc) HandlerFunc {
				return middleware(r.Middleware(next))
			}
		} else if middleware != nil {
			m = middleware
		} else if r.Middleware != nil {
			m = r.Middleware
		}

		if len(r.Children) > 0 {
			cn.generateFromRoutes(r.Children, route, m)
		} else {
			if m != nil {
				cn.handlers[r.Method] = m(r.Handler)
			} else {
				cn.handlers[r.Method] = r.Handler
			}
			cn.route = route
		}
	}
}

func (n *node) generateFromPath(path string) *node {
	if i := strings.IndexByte(path, '/'); i > -1 {
		return n.generateNode(path[:i]).generateFromPath(path[i+1:])
	} else {
		return n.generateNode(path)
	}
}

func (n *node) generateNode(path string) *node {
	if path == "" {
		return n
	} else if path == "*" {
		n.wildcard = true
		return n
	} else {
		return n.findOrCreateNode(path)
	}
}

func (n *node) findOrCreateNode(path string) *node {
	k := path
	if path[0] == ':' {
		k = paramKey
	}
	if n.children[k] == nil {
		n.children[k] = newNode(n)
	}
	return n.children[k]
}

func (n *node) findNode(path string) *node {
	var node *node

	if i := strings.IndexByte(path, '/'); i > -1 {
		cn := n.findPart(path[:i])
		if cn != nil {
			node = cn.findNode(path[i+1:])
		} else {
			node = nil
		}
	} else {
		node = n.findPart(path)
	}

	if node == nil {
		return n.findWildcardNode()
	} else if len(node.handlers) == 0 {
		return node.findWildcardNode()
	} else {
		return node
	}
}

func (n *node) findPart(part string) *node {
	if part == "" {
		return n
	} else if cn, ok := n.children[part]; ok {
		return cn
	} else if cn, ok := n.children[paramKey]; ok {
		return cn
	} else {
		return nil
	}
}

func (n *node) findWildcardNode() *node {
	if n.wildcard {
		return n
	}
	if n.parent == nil {
		return nil
	}
	return n.parent.findWildcardNode()
}
