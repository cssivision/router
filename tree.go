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

func (n *node) insert(method, pattern string, handler Handle, ignore bool) {
	if ignore {
		pattern = strings.ToLower(pattern)
	}
	pattern = strings.TrimPrefix(pattern, "/")

	p := n
	frags := strings.Split(pattern, "/")
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
				panic(fmt.Sprintf(`invalid named parameter: "%s"`, name))
			}
			nn.name = name

			if frag[0] == '*' {
				nn.wildcard = true
			}

			if child := p.parameterChild; child != nil {
				if child.name != name || child.wildcard != nn.wildcard {
					panic(fmt.Sprintf(`invalid named parameter: "%s"`, name))
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

	if p.handlers[method] != nil {
		panic("conflicts with existing " + pattern + ", method " + method)
	}

	p.handlers[method] = handler
	p.pattern = pattern
}

func printTree(root *node) {
	if root == nil {
		return
	}

	if len(root.handlers) > 0 {
		fmt.Println("pattern: ", root.pattern)
	}

	for _, v := range root.children {
		printTree(v)
	}

	if root.parameterChild != nil {
		printTree(root.parameterChild)
	}
}

func (n *node) find(path, method string, ignoreCase bool, trailingSlashRedirect bool) (Handle, Params, bool) {
	if path == "" || path[0] != '/' {
		panic(fmt.Errorf(`path must start with "/": "%s"`, path))
	}

	if ignoreCase {
		path = strings.ToLower(path)
	}

	path = strings.TrimPrefix(path, "/")
	var tsr bool
	var matchedParams map[string]string
	p := n
	frags := strings.Split(path, "/")
	for index, frag := range frags {
		nn := p.children[frag]
		if nn == nil {
			nn = p.parameterChild
		}

		if nn == nil {
			// TrailingSlashRedirect: /a/b/ -> /a/b
			if trailingSlashRedirect && p.endpoint && index == len(frags)-1 && frag == "" {
				tsr = true
			}
			return nil, matchedParams, tsr
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

	if trailingSlashRedirect && p.children[""] != nil {
		// TrailingSlashRedirect: /a/b -> /a/b/
		tsr = true
	}

	return p.handlers[method], matchedParams, tsr
}
