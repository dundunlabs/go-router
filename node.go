package gorouter

import (
	"strings"
)

func newNode(route string, parent *node) *node {
	return &node{
		route:    route,
		parent:   parent,
		children: make(map[string]*node),
		handlers: make(map[string]Handler),
	}
}

type node struct {
	route    string
	parent   *node
	children map[string]*node
	handlers map[string]Handler
}

func (n *node) generateFromRoutes(routes []Route) {
	for _, r := range routes {
		cn := n.generateFromPath(r.Path)
		if len(r.Children) > 0 {
			cn.generateFromRoutes(r.Children)
		} else {
			cn.handlers[r.Method] = r.Handler
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
	} else {
		return n.findOrCreateNode(path)
	}
}

func (n *node) findOrCreateNode(path string) *node {
	k := path
	if path[0] == ':' {
		k = "#"
	}
	if n.children[k] == nil {
		n.children[k] = newNode(n.route+"/"+path, n)
	}
	return n.children[k]
}

func (n *node) findNode(path string) *node {
	if i := strings.IndexByte(path, '/'); i > -1 {
		cn := n.findPart(path[:i])
		if cn != nil {
			return cn.findNode(path[i+1:])
		}
		return nil
	} else {
		return n.findPart(path)
	}
}

func (n *node) findPart(part string) *node {
	if part == "" {
		return n
	} else if cn, ok := n.children[part]; ok {
		return cn
	} else if cn, ok := n.children["#"]; ok {
		return cn
	} else {
		return nil
	}
}
