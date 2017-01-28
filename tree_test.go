package router

import (
    "testing"
)

func TestInsert(t *testing.T) {
    t.Run("test for path /", func(t *testing.T) {
        n := &node{
            children: make(map[string]*node),
            handlers: make(map[string]Handle),
        }
    })
}