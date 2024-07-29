package main

import (
	"bufio"
	"fmt"
)

// Let's allow potential word overlap between two consecutive
// calls to Write, e.g.
//	c.Write("hel"); c. Write("lo ") should count only one word, "hello".
type WordCounter struct {
	n int
	p []byte
}

type LineCounter struct {
	n int
	p []byte
}

// « The Go compiler does not support accessing a struct field
// x.f where x is of type parameter type even if all types in
// the type parameter's type set have a field f. We may remove
// this restriction in a future release.
//
// Generics aren't covered by gopl anyway; this also means
// we can't use a generic doWrite that access our Counter's fields.
func (c *WordCounter) Incr() { c.n += 1 }
func (c *LineCounter) Incr() { c.n += 1 }

func (c *WordCounter) Save(p []byte) { c.p = p }
func (c *LineCounter) Save(p []byte) { c.p = p }

func (c *WordCounter) GetSave() []byte { return c.p }
func (c *LineCounter) GetSave() []byte { return c.p }

type Counter interface {
	Incr()
	Save(p []byte)
	GetSave() []byte
}

// Perhaps there's a more efficient approach, but this seems to work.
func doWrite(c Counter, p []byte, f bufio.SplitFunc) (int, error) {
	// Grab what remained last time.
	q := append(c.GetSave(), p...)

	m := 0
	for {
		n, ts, err := f(q[m:], false)
//		fmt.Printf("n=%d; ts=%s; err=%v\n", n, string(ts), err)
		if n == 0 {
			break
		}
		if err != nil {
			return m, err
		}
		m += n

		// we may be advancing further, but still not reaching
		// a word (spaces)
		if len(ts) > 0 {
			c.Incr()
		}
	}
	c.Save(q[m:])
//	println("Saving for later: '"+string(c.GetSave())+"'")
	return m, nil
}

func (c *WordCounter) Write(p []byte) (int, error) {
	return doWrite(c, p, bufio.ScanWords)
}

func (c *LineCounter) Write(p []byte) (int, error) {
	return doWrite(c, p, bufio.ScanLines)
}

func main() {
	c := &WordCounter{0,make([]byte,0)}

	c.Write([]byte("hello, world"))

	fmt.Printf("Counter: %d\n", c.n)
	c.Write([]byte(" "))
	fmt.Printf("Counter: %d\n", c.n)
}
