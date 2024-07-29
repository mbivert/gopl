package main

import (
	"bytes"
	"fmt"
)

type tree struct {
	value int
	left, right *tree
}
func add(t *tree, value int) *tree {
	if t == nil {
		return &tree{value: value}
	}
	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}
	return t
}

func addAll(t *tree, values ...int) *tree {
	for _, value := range values {
		t = add(t, value)
	}
	return t
}

// « Write a String method for the *tree type in
// gopl.io/ch4/treesort (§4.4) reveals the sequence
// of values in the tre e. »
func (t *tree) String() string {
	var b bytes.Buffer

	if t == nil {
		return ""
	}

	// != nil guards only for space display control then.
	if t.left != nil {
		b.WriteString(t.left.String())
		b.WriteByte(' ')
	}
	fmt.Fprintf(&b, "%d", t.value)
	if t.right != nil {
		b.WriteByte(' ')
		b.WriteString(t.right.String())
	}

	return b.String()
}
