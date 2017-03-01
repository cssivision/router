package router

import (
	"net/http"
	"strings"
)

var default405Body = []byte("405 method not allowed")

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
	NoRoute http.Handler

	// Configurable http.Handler which is called when method is not allowed. If it is not set, http.NotFound is used.
	NoMethod http.Handler

	// Methods which has been registered
	allowMethods map[string]bool
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
			basePath: "",
		},
		tree: &node{
			children: make(map[string]*node),
			handlers: make(map[string]Handle),
		},
		TrailingSlashRedirect: true,
		allowMethods:          make(map[string]bool),
	}

	router.RouterPrefix.router = router

	return router
}

// ServeHTTP makes the router implement the http.Handler interface.
func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if !r.allowMethods[req.Method] {
		if r.NoMethod != nil {
			r.NoMethod.ServeHTTP(rw, req)
		} else {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			rw.Write(default405Body)
		}
		return
	}

	var pattern string
	pattern = req.URL.Path
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
		} else if len(path) == 1 {
			// do nothing
		} else {
			req.URL.Path = path + "/"
		}

		if tsr {
			pattern = req.URL.Path
			http.Redirect(rw, req, pattern, http.StatusMovedPermanently)
			return
		}
	}

	if r.NoRoute != nil {
		r.NoRoute.ServeHTTP(rw, req)
	} else {
		http.NotFound(rw, req)
	}
}
