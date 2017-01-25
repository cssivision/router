package router

import (
	"strings"
)

type node struct {
	children map[string]*node
	handlers map[string]Handle
}

func (n *node) insert(method, pattern string, handler Handle, ignore bool) {
	if ignore {
		pattern = strings.ToLower(pattern)
	}

	pattern = strings.Trim(pattern, "/")
	p := n

	frags := strings.Split(pattern, "/")
	for _, frag := range frags {
		if p.children == nil {
			p.children = make(map[string]*node)
		}

		if p.children[frag] == nil {
			p.children[frag] = new(node)
		}

		p = p.children[frag]
	}

	if p.handlers == nil {
		p.handlers = make(map[string]Handle)
	}

	if p.handlers[method] != nil {
		panic("conflicts with existing " + pattern + " method " + method)
	}

	p.handlers[method] = handler
}

func (n *node) find(path, method string, ignore bool) (Handle, Params) {
	if ignore {
		path = strings.ToLower(path)
	}

	path = strings.TrimPrefix(path, "/")
	var frags []string
	p := n
	frags = strings.Split(path, "/")
	for _, frag := range frags {
		if p.children == nil {
			return nil, Params{}
		}

		if p.children[frag] == nil {
			return nil, Params{}
		}

		p = p.children[frag]
	}

	if p.handlers == nil {
		return nil, Params{}
	}

	if p.handlers[method] == nil {
		return nil, Params{}
	}

	return p.handlers[method], Params{}
}
