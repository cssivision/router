package router

import (
	"net/http"
)

type Router struct {
	tree *node
}

type Handle func(http.ResponseWriter, *http.Request, Params)
type Param struct {
	Key, Value string
}

type Params []Param

func New() *Router {
	return &Router{
		tree: new(node),
	}
}

func (r *Router) Get(pattern string, handle Handle) {
	r.Handle(http.MethodGet, pattern, handle)
}

func (r *Router) Post(pattern string, handle Handle) {
	r.Handle(http.MethodPost, pattern, handle)
}

func (r *Router) Put(pattern string, handle Handle) {
	r.Handle(http.MethodPut, pattern, handle)
}

func (r *Router) Handle(method, pattern string, handle Handle) {
	if pattern[0] != '/' {
		panic("path must begin with '/' in path '" + pattern + "'")
	}

	r.tree.insert(method, pattern, handle)
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	handle, _ := r.tree.find(req.URL.String(), req.Method)

	if handle == nil {
		http.NotFound(rw, req)
		return
	}

	handle(rw, req, Params{})
}
