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
type Param struct {
	Key, Value string
}

type Params []Param

func New() *Router {
	return &Router{
		tree:                  new(node),
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

func (r *Router) Handle(method, pattern string, handler Handle) {
	if pattern[0] != '/' {
		panic("path must begin with '/' in path '" + pattern + "'")
	}

	if r.tree == nil {
		r.tree = new(node)
	}

	r.tree.insert(method, pattern, handler, r.IgnoreCase)
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	handle, _ := r.tree.find(req.URL.String(), req.Method, r.IgnoreCase)

	if handle != nil {
		handle(rw, req, Params{})
		return
	}

	path := req.URL.Path

	if r.TrailingSlashRedirect {
		if len(path) > 1 && path[len(path)-1] == '/' {
			req.URL.Path = path[:len(path)-1]
			handle, _ = r.tree.find(req.URL.String(), req.Method, r.IgnoreCase)
		} else {
			req.URL.Path = path + "/"
			handle, _ = r.tree.find(req.URL.String(), req.Method, r.IgnoreCase)
		}

		if handle != nil {
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
