package main

// NOTE: I'm not sure whether I understood the exercise correctly:
// what is implemented below is close to chapter 5's fetch, but the URLs
// provided  on the command lines are fetched in parallel, and only the
// fastest one is fetched, while the others are cancelled.
//
// This doesn't feel really useful of a program, hence why I'm not sure
// as to whether I've understood the exercise.
//
// NOTE: we avoid fetching content to disk to avoid having to deal with
// distinct filenames.

// go run fetch.go https://tales.mbivert.com/ https://google.fr/ https://mbivert.com/
//
// 	-> Last one is the smallest, and thus should almost always be the fastest.

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var done = make(chan struct{})

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
	}
	return false
}

type myCtx struct{}

func (c myCtx) Deadline() (time.Time, bool) { return time.Time{}, false }
func (c myCtx) Done() <-chan struct{}       { return done }
func (c myCtx) Err() error {
	// NOTE: those rarely happen, but can still be observed once in a while.
	if cancelled() {
		return fmt.Errorf("cancelled")
	}
	return nil
}
func (c myCtx) Value(key interface{}) interface{} { return nil }

// Fetch downloads the URL and returns the
// name and length of the local file.
func fetch(url string) (content string, n int64, err error) {
	req, err := http.NewRequestWithContext(myCtx{}, "GET", url, nil)
	if err != nil {
		return "", 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}

	defer resp.Body.Close()
	var buf strings.Builder
	n, err = io.Copy(&buf, resp.Body)
	return buf.String(), n, err
}

type ret struct {
	url     string
	content string
	n       int64
}

func fetchFastest(urls []string) (string, string, int64) {
	responses := make(chan ret, len(urls))
	for _, url := range urls {
		go func() {
			content, n, err := fetch(url)
			if err != nil {
				fmt.Fprintf(os.Stderr, "fetch %s: %v\n", url, err)
			}
			responses <- ret{url, content, n}
		}()
	}
	x := <-responses
	close(done)
	return x.url, x.content, x.n
}

func main() {
	url, content, n := fetchFastest(os.Args[1:])
	fmt.Printf("%s (%d bytes).\n\n", url, n)
	fmt.Printf("%s", content)
}
