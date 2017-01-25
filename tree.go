package router

import (
    "strings"
)

type node struct {
    pattern  string
    children map[string]*node
    handle   map[string]Handle
}

func (n *node) insert(method, pattern string, handle Handle) {
    pattern = strings.Trim(pattern, "/")
    var frags []string
    p := n
    if len(pattern) == 0 {
        frags = []string{}
    } else {
       frags = strings.Split(pattern, "/")
    }

    for _, frag := range frags {
        if p.children == nil {
            p.children = make(map[string]*node)
        }

        if p.children[frag] == nil {
            p.children[frag] = &node{
                pattern: frag,
            }
        }

        p = p.children[frag]
    }

    if p.handle == nil {
        p.handle = make(map[string]Handle)
    }

    if p.handle[method] != nil {
        panic("handle exist")
    }

    p.handle[method] = handle
}

func (n *node) find(path, method string) (handle Handle, ps Params) {
    path = strings.Trim(path, "/")
    var frags []string
    p := n
    if len(path) == 0 {
        frags = []string{}
    } else {
       frags = strings.Split(path, "/")
    }

    for _, frag := range frags {
        if p.children == nil {
            return nil, Params{}
        }

        if p.children[frag] == nil {
            return nil, Params{}
        }

        p = p.children[frag]
    }

    if p.handle == nil {
        return nil, Params{}
    }

    if p.handle[method] == nil {
        return nil, Params{}
    }

    return p.handle[method], Params{}
}