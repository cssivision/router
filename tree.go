package router

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	nameRegexp = regexp.MustCompile(`^\w+$`)
)

type node struct {
	pattern        string
	name           string
	endpoint       bool
	wildcard       bool
	parameterChild *node
	children       map[string]*node
	handlers       map[string]Handle
}

func (n *node) insert(pattern string) *node {
	if strings.Contains(pattern, "//") {
		panic(fmt.Errorf(`must not contain multi-slash: "%s"`, pattern))
	}

	pattern = strings.TrimPrefix(pattern, "/")
	frags := strings.Split(pattern, "/")

	p := n
	for index, frag := range frags {

		if p.children[frag] != nil {
			p = p.children[frag]
			continue
		}

		nn := &node{
			children: make(map[string]*node),
			handlers: make(map[string]Handle),
		}

		if frag == "" {
			p.children[frag] = nn
		} else if frag[0] == '*' || frag[0] == ':' {
			name := frag[1:]
			if !nameRegexp.MatchString(name) {
				panic(fmt.Sprintf(`invalid named parameter: "%s"`, pattern))
			}
			nn.name = name

			if frag[0] == '*' {
				nn.wildcard = true
			}

			if child := p.parameterChild; child != nil {
				if child.name != name || child.wildcard != nn.wildcard {
					panic(fmt.Sprintf(`invalid named parameter: "%s"`, pattern))
				}
				p = child
				continue
			} else {
				p.parameterChild = nn
			}
		} else {
			p.children[frag] = nn
		}

		p = nn

		if index == len(frags)-1 {
			nn.endpoint = true
			continue
		}

		if nn.wildcard {
			panic("can't define path after wildcard pattern")
		}
	}

	p.pattern = "/" + pattern
	return p
}

func (n *node) addHandle(method string, handler Handle) {
	if n.handlers[method] != nil {
		panic("pattern: " + n.pattern + ", method: " + method + ", conflicts with existing route")
	}

	n.handlers[method] = handler
}

func (n *node) find(path string) (*node, Params, bool) {
	if path == "" || path[0] != '/' {
		panic(fmt.Errorf(`path must start with "/": "%s"`, path))
	}

	var tsr bool
	var matchedParams map[string]string
	path = strings.TrimPrefix(path, "/")
	frags := strings.Split(path, "/")

	p := n
	for index, frag := range frags {
		nn := p.children[frag]
		if nn == nil {
			nn = p.parameterChild
		}

		if nn == nil {
			// TrailingSlashRedirect: /a/b/ -> /a/b
			if p.endpoint && index == len(frags)-1 && frag == "" {
				tsr = true
			}
			return p, matchedParams, tsr
		}

		p = nn

		if p.name != "" {
			if matchedParams == nil {
				matchedParams = make(map[string]string)
			}

			if p.wildcard {
				fmt.Println(strings.Join(frags[index:], "/"))
				matchedParams[p.name] = strings.Join(frags[index:], "/")
				break
			} else {
				matchedParams[p.name] = frag
			}
		}
	}

	if p.children[""] != nil {
		// TrailingSlashRedirect: /a/b -> /a/b/
		tsr = true
	}

	return p, matchedParams, tsr
}
