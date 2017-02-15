package router

import (
	"net/http"
	"strings"
)

// Router is a http.Handler which can be used to dispatch requests to different
// handler functions via configurable routes
type Router struct {
	RouterPrefix

	// tree used to keep handler with path
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
	router := &Router{
		RouterPrefix: RouterPrefix{
			BasePath: "",
		},
		tree: &node{
			children: make(map[string]*node),
			handlers: make(map[string]Handle),
		},
		TrailingSlashRedirect: true,
	}

	router.RouterPrefix.Router = router

	return router
}

// ServeHTTP makes the router implement the http.Handler interface.
func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var pattern string
	pattern = req.URL.String()
	if r.IgnoreCase {
		pattern = strings.ToLower(pattern)
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
