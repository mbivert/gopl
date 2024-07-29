package main

import (
	"testing"
	"fmt"
)

func TestWordCounter(t *testing.T) {

	c := &WordCounter{0,make([]byte,0)}

	c.Write([]byte("hello, world"))

	if c.n != 1 {
		t.Errorf("'hello,' should have been counted")
	}

	c.Write([]byte(" "))

	if c.n != 2 {
		t.Errorf("'world' should have been counted")
	}

	c.Write([]byte("o"))

	if c.n != 2 {
		t.Errorf("not yet a word")
	}

	c.Write([]byte("f"))

	if c.n != 2 {
		t.Errorf("not yet a word (bis)")
	}

	// NOTE: will call ScanWords twice: a first time because
	// we consume "of", and a second time to read that second
	// space (and generate no words).
	c.Write([]byte("  "))

	if c.n != 3 {
		t.Errorf("'of' should have been read")
	}
}

func TestLineCounter(t *testing.T) {

	c := &LineCounter{0,make([]byte,0)}

	c.Write([]byte("hello, world"))

	if c.n != 0 {
		t.Errorf("Not yet a line")
	}

	c.Write([]byte(" "))

	if c.n != 0 {
		t.Errorf("Not yet a line")
	}


	c.Write([]byte("o"))

	if c.n != 0 {
		t.Errorf("Not yet a line")
	}


	c.Write([]byte("f"))

	if c.n != 0 {
		t.Errorf("Not yet a line")
	}

	c.Write([]byte("  \n"))

	if c.n != 1 {
		t.Errorf("One first line")
	}
}

func TestCountingWriter(t *testing.T) {
	c, n := CountingWriter(&WordCounter{0,make([]byte,0)})

	// Haven't tested this before here actually (used as an io.Writer)
	fmt.Fprintf(c, "hello, world")

	if *n != int64(len("hello, ")) {
		t.Errorf("We should have counted %d bytes (`hello, `); have %d",
			len("hello "), *n)
	}

	c.Write([]byte(" "))

	if *n != int64(len("hello, world ")) {
		t.Errorf("We should have counted %d bytes (`hello, `); have %d",
			len("hello, world "), *n)
	}
}
