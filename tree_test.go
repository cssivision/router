package router

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsert(t *testing.T) {
	t.Run("test for path /", func(t *testing.T) {
		tree1 := New().tree
		tree2 := New().tree
		n1 := tree1.insert("/")
		n2 := tree2.insert("")

		assert.Equal(t, n1.name, "", fmt.Sprintf("got node name %s, expected %s", tree1.name, ""))
		assert.Equal(t, n2.name, "", fmt.Sprintf("got node name %s, expected %s", tree2.name, ""))
		assert.Equal(t, n1, tree1.insert("/"), "insert same pattern, should return same tree node")
		assert.Equal(t, n1, tree1.insert(""), "insert same pattern, should return same tree node")
		assert.Equal(t, n2, tree2.insert("/"), "insert same pattern, should return same tree node")
		assert.Equal(t, n2, tree2.insert(""), "insert same pattern, should return same tree node")
	})

	t.Run("test for simple path", func(t *testing.T) {
		tree := New().tree
		n := tree.insert("/a/b")

		assert.Equal(t, n.name, "", fmt.Sprintf("got node name %s, expected %s", n.name, ""))
		assert.Equal(t, n, tree.insert("/a/b"), "same pattern, should return same tree node")
		assert.Equal(t, n, tree.insert("a/b"), "same pattern, should return same tree node")
		assert.NotEqual(t, n, tree.insert("/a/b/"), "different pattern, should return different tree node")
		assert.NotEqual(t, n, tree.insert("a/b/"), "different pattern, should return different tree node")
		assert.Panics(t, func() {
			tree.insert("/a//b")
		}, fmt.Sprintf(`must not contain multi-slash: "%s"`, "/a//b"))
	})

	t.Run("test for named pattern", func(t *testing.T) {
		tree := New().tree
		n := tree.insert("/a/:b")
		matched, ps, _ := tree.find("/a/name")

		assert.Equal(t, ps["b"], "name", fmt.Sprintf("got params b: %s, expected %s", ps["b"], "name"))
		assert.Equal(t, n, matched, fmt.Sprintf("same pattern, should return same tree node"))
		assert.Equal(t, matched.name, "b", fmt.Sprintf("got params name: %s, expected %s", matched.name, "b"))
		assert.Panics(t, func() {
			tree.insert("/:$~!")
		})
		assert.Nil(t, matched.parameterChild, "should have parameterChild")
		assert.Equal(t, matched.pattern, "/a/:b")
		assert.Panics(t, func() {
			tree.insert("/a/:x")
		})

		assert.Panics(t, func() {
			tree.insert("/a/:b/c")
			tree.insert("/a/:x/c")
		})
	})

	t.Run("test for wildcard pattern", func(t *testing.T) {
		tree := New().tree

		assert.Panics(t, func() {
			tree.insert("/a/*")
		})

		n := tree.insert("/a/*b")
		assert.Equal(t, n, tree.insert("/a/*b"))
		assert.Equal(t, n, tree.insert("a/*b"))
		assert.Equal(t, n.name, "b")
		assert.True(t, n.wildcard)
		assert.Equal(t, n.pattern, "/a/*b")
		assert.Panics(t, func() {
			tree.insert("/a/*c")
		})
	})
}

func TestFind(t *testing.T) {

}
