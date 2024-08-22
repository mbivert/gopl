package main

import (
	"testing"
	// "fmt"
)

func TestTreeString(t *testing.T) {
	var p *tree

	p = addAll(p, 10, 5, 7, 1, 42, 12)

	if p.String() != "1 5 7 10 12 42" {
		t.Errorf("Unexpected tree string: %s", p.String())
	}
}
