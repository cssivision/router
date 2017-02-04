package router

import (
	"net/http"
	"strings"
)

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
type Router struct {
	tree *node

	// Ignore case when matching URL path.
	IgnoreCase bool

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// TrailingSlashRedirect: /a/b/ -> /a/b
	// TrailingSlashRedirect: /a/b -> /a/b/
	TrailingSlashRedirect bool

	// Configurable http.Handler which is called when no matching route is
	// found. If it is not set, http.NotFound is used.
	NotFound http.Handler
}

// Handle is a function that can be registered to a route to handle HTTP
// requests. Like http.HandlerFunc, but has a third parameter for the
// values of named/wildcards parameters.
type Handle func(http.ResponseWriter, *http.Request, Params)

// Param is a single URL parameter, a map[string]string.
type Params map[string]string

// New returns a new initialized Router, with default configuration
func New() *Router {
	return &Router{
		tree: &node{
			children: make(map[string]*node),
			handlers: make(map[string]Handle),
		},
		IgnoreCase:            true,
		TrailingSlashRedirect: true,
	}
}

// Get is a shortcut for router.Handle("GET", path, handle)
func (r *Router) Get(pattern string, handler Handle) {
	r.Handle(http.MethodGet, pattern, handler)
}

// Post is a shortcut for router.Handle("POST", path, handle)
func (r *Router) Post(pattern string, handler Handle) {
	r.Handle(http.MethodPost, pattern, handler)
}

// Put is a shortcut for router.Handle("PUT", path, handle)
func (r *Router) Put(pattern string, handler Handle) {
	r.Handle(http.MethodPut, pattern, handler)
}

// Delete is a shortcut for router.Handle("DELETE", path, handle)
func (r *Router) Delete(pattern string, handler Handle) {
	r.Handle(http.MethodDelete, pattern, handler)
}

// Options is a shortcut for router.Handle("OPTIONS", path, handle)
func (r *Router) Options(pattern string, handler Handle) {
	r.Handle(http.MethodOptions, pattern, handler)
}

// Trace is a shortcut for router.Handle("TRACE", path, handle)
func (r *Router) Trace(pattern string, handler Handle) {
	r.Handle(http.MethodTrace, pattern, handler)
}

// Head is a shortcut for router.Handle("HEAD", path, handle)
func (r *Router) Head(pattern string, handler Handle) {
	r.Handle(http.MethodHead, pattern, handler)
}

// Patch is a shortcut for router.Handle("PATCH", path, handle)
func (r *Router) Patch(pattern string, handler Handle) {
	r.Handle(http.MethodPatch, pattern, handler)
}

// Add prefix for a router, and return a new one
func (r *Router) Prefix(prefix string) *RouterPrefix {
	if prefix == "" {
		panic("prefix can not be empty")
	}

	if prefix[0] != '/' {
		panic("prefix must begin with /")
	}

	prefix = strings.TrimSuffix(prefix, "/")

	return &RouterPrefix{
		BasePath: prefix,
		Router:   r,
	}
}

// Handle registers a new request handle with the given path and method.
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
func (r *Router) Handle(method, pattern string, handler Handle) {
	if method == "" {
		panic("invalid http method")
	}

	if pattern[0] != '/' {
		panic("path must begin with '/', '" + pattern + "'")
	}

	if r.tree == nil {
		r.tree = &node{
			children: make(map[string]*node),
			handlers: make(map[string]Handle),
		}
	}

	if r.IgnoreCase {
		pattern = strings.ToLower(pattern)
	}
	r.tree.insert(pattern).addHandle(method, handler)
}

// ServeHTTP makes the router implement the http.Handler interface.
func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var pattern string
	if r.IgnoreCase {
		pattern = strings.ToLower(req.URL.String())
	}

	n, ps, tsr := r.tree.find(pattern)
	var handle Handle
	if n != nil {
		handle = n.handlers[req.Method]
	}

	if handle != nil {
		handle(rw, req, ps)
		return
	}

	path := req.URL.Path
	if r.TrailingSlashRedirect {
		if len(path) > 1 && path[len(path)-1] == '/' {
			req.URL.Path = path[:len(path)-1]
		} else {
			req.URL.Path = path + "/"
		}

		if tsr {
			pattern = strings.ToLower(req.URL.String())
			http.Redirect(rw, req, pattern, http.StatusMovedPermanently)
			return
		}
	}

	if r.NotFound != nil {
		r.NotFound.ServeHTTP(rw, req)
	} else {
		http.NotFound(rw, req)
	}
}
