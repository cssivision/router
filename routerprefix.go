package router

import (
	"net/http"
	"strings"
)

// RouterPrefix
type RouterPrefix struct {
	Router *Router

	// Prefix path of a router
	BasePath string
}

// Get is a shortcut for router.Handle("GET", path, handle) with BasePath
func (r *RouterPrefix) Get(pattern string, handler Handle) {
	r.Handle(http.MethodGet, pattern, handler)
}

// Post is a shortcut for router.Handle("POST", path, handle) with BasePath
func (r *RouterPrefix) Post(pattern string, handler Handle) {
	r.Handle(http.MethodPost, pattern, handler)
}

// Put is a shortcut for router.Handle("PUT", path, handle) with BasePath
func (r *RouterPrefix) Put(pattern string, handler Handle) {
	r.Handle(http.MethodPut, pattern, handler)
}

// Delete is a shortcut for router.Handle("DELETE", path, handle) with BasePath
func (r *RouterPrefix) Delete(pattern string, handler Handle) {
	r.Handle(http.MethodDelete, pattern, handler)
}

// Options is a shortcut for router.Handle("OPTIONS", path, handle) with BasePath
func (r *RouterPrefix) Options(pattern string, handler Handle) {
	r.Handle(http.MethodOptions, pattern, handler)
}

// Trace is a shortcut for router.Handle("TRACE", path, handle) with BasePath
func (r *RouterPrefix) Trace(pattern string, handler Handle) {
	r.Handle(http.MethodTrace, pattern, handler)
}

// Head is a shortcut for router.Handle("HEAD", path, handle) with BasePath
func (r *RouterPrefix) Head(pattern string, handler Handle) {
	r.Handle(http.MethodHead, pattern, handler)
}

// Patch is a shortcut for router.Handle("PATCH", path, handle) with BasePath
func (r *RouterPrefix) Patch(pattern string, handler Handle) {
	r.Handle(http.MethodPatch, pattern, handler)
}

// Add prefix for a router, and return a new one with BasePath
func (r *RouterPrefix) Prefix(prefix string) *RouterPrefix {

	return &RouterPrefix{
		BasePath: prefix,
		Router:   r.Router,
	}
}

// Handle registers a new request handle with the given path and method.
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
func (r *RouterPrefix) Handle(method, pattern string, handler Handle) {
	if r.BasePath != "" {
		pattern = r.BasePath + pattern
	}

	// r.Router.Handle(method, pattern, handler)

	if method == "" {
		panic("invalid http method")
	}

	if pattern[0] != '/' {
		panic("path must begin with '/', '" + pattern + "'")
	}

	if r.Router.tree == nil {
		r.Router.tree = &node{
			children: make(map[string]*node),
			handlers: make(map[string]Handle),
		}
	}

	if r.Router.IgnoreCase {
		pattern = strings.ToLower(pattern)
	}
	r.Router.tree.insert(pattern).addHandle(method, handler)
}
