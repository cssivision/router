package router

import (
	"net/http"
)

type Router struct {
	tree                  *node
	IgnoreCase            bool
	TrailingSlashRedirect bool
	NotFound              http.Handler
}

type Handle func(http.ResponseWriter, *http.Request, Params)

type Params map[string]string

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

func (r *Router) Get(pattern string, handler Handle) {
	r.Handle(http.MethodGet, pattern, handler)
}

func (r *Router) Post(pattern string, handler Handle) {
	r.Handle(http.MethodPost, pattern, handler)
}

func (r *Router) Put(pattern string, handler Handle) {
	r.Handle(http.MethodPut, pattern, handler)
}

func (r *Router) Delete(pattern string, handler Handle) {
	r.Handle(http.MethodDelete, pattern, handler)
}

func (r *Router) Options(pattern string, handler Handle) {
	r.Handle(http.MethodOptions, pattern, handler)
}

func (r *Router) Trace(pattern string, handler Handle) {
	r.Handle(http.MethodTrace, pattern, handler)
}

func (r *Router) Head(pattern string, handler Handle) {
    r.Handle(http.MethodHead, pattern, handler)
}

func (r *Router) Patch(pattern string, handler Handle) {
    r.Handle(http.MethodPatch, pattern, handler)
}

func (r *Router) Handle(method, pattern string, handler Handle) {
	if pattern[0] != '/' {
		panic("path must begin with '/' in path '" + pattern + "'")
	}

	if r.tree == nil {
		r.tree = &node{
			children: make(map[string]*node),
			handlers: make(map[string]Handle),
		}
	}

	r.tree.insert(method, pattern, handler, r.IgnoreCase)
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	handle, ps, tsr := r.tree.find(req.URL.String(), req.Method, r.IgnoreCase, r.TrailingSlashRedirect)

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
			http.Redirect(rw, req, req.URL.String(), http.StatusMovedPermanently)
			return
		}
	}

	if r.NotFound != nil {
		r.NotFound.ServeHTTP(rw, req)
	} else {
		http.NotFound(rw, req)
	}
}
