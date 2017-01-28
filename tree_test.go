package router

import (
	"fmt"
	"testing"
)

func TestInsert(t *testing.T) {
	t.Run("test for path /", func(t *testing.T) {
		tree1 := New().tree
		tree2 := New().tree
		n1 := tree1.insert("/")
		n2 := tree2.insert("")

		if n1.name != "" {
			t.Errorf("got node name %s, expected %s", tree1.name, "")
		}

		if n2.name != "" {
			t.Errorf("got node name %s, expected %s", tree2.name, "")
		}

		if n1 != tree1.insert("/") || n1 != tree1.insert("") {
			t.Errorf("insert same pattern, should return same tree node")
		}

		if n2 != tree2.insert("/") || n2 != tree2.insert("") {
			t.Errorf("insert same pattern, should return same tree node")
		}
	})

	t.Run("test for simple path", func(t *testing.T) {
		tree := New().tree
		n := tree.insert("/a/b")

		if n.name != "" {
			t.Errorf("got node name %s, expected %s", n.name, "")
		}

		if n != tree.insert("/a/b") {
			t.Errorf("same pattern, should return same tree node")
		}

		if n != tree.insert("a/b") {
			t.Errorf("same pattern, should return same tree node")
		}

		if n == tree.insert("/a/b/") {
			t.Errorf("different pattern, should return different tree node")
		}

		if n == tree.insert("a/b/") {
			t.Errorf("different pattern, should return different tree node")
		}

		defer func() {
			if rec := recover(); rec != nil {
				if rec.(error).Error() != fmt.Sprintf(`must not contain multi-slash: "%s"`, "/a//b") {
					t.Errorf(rec.(error).Error())
				}
			}
		}()

		tree.insert("/a//b")
	})

	t.Run("test for named pattern", func(t *testing.T) {

	})
}
