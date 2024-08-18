package main

// NOTE: cancelling the HTTP requests is merely useful to speed up the
// shutdown; we still need to wire everything else properly.
//
// NOTE: the authors advise to use Request.Cancel (<- chan struct{}),
// but the field is now deprecated in favors of context.

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/net/html"
)

var done = make(chan struct{})

var depth int
var ncrawl int

func cancelled() bool {
	// NOTE: we can't "case _, ok <- done" without
	// a default, because it would block if done isn't
	// closed
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

// As stated earlier in this file, request.Cancel is deprecated in
// favor of contexts; we must thus create a context wrapping access
// to our "done" channel. The error is necessary (panic() otherwise)
type myCtx struct{}

func (c myCtx) Deadline() (time.Time, bool) { return time.Time{}, false }
func (c myCtx) Done() <-chan struct{}       { return done }
func (c myCtx) Err() error {
	if cancelled() {
		return fmt.Errorf("cancelled")
	}
	return nil
}
func (c myCtx) Value(key interface{}) interface{} { return nil }

func extractLinks(url string) ([]string, error) {
	req, err := http.NewRequestWithContext(myCtx{}, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

type Link struct {
	depth int
	url   string
}

func mkLinks(depth int, urls []string) []Link {
	links := []Link{}
	for _, url := range urls {
		links = append(links, Link{depth, url})
	}
	return links
}

func crawl(link Link) []Link {
	fmt.Println(link.url)

	if link.depth == depth {
		return []Link{}
	}
	urls, err := extractLinks(link.url)
	if err != nil {
		log.Print(err)
	}
	return mkLinks(link.depth+1, urls)
}

// NOTE: this is using the "alternative" version of the crawler, extended
// to (hopefully) fix the termination problem.
func main() {
	flag.IntVar(&depth, "depth", 3, "maximum depth")
	flag.IntVar(&ncrawl, "ncrawl", 20, "number of concurrent crawlers")

	flag.Parse()

	var n int                      // number of pending links to unseenLinks
	worklist := make(chan []Link)  // lists of URLs, may have duplicates
	unseenLinks := make(chan Link) // de-duplicated URLs

	// Add command-line arguments to worklist.
	n = len(flag.Args())

	// perhaps we'd want to check for done here too, thought it's less
	// likely to be an issue, it still is theoretically
	go func() { worklist <- mkLinks(0, flag.Args()) }()

	go func() {
		os.Stdin.Read(make([]byte, 1))
		fmt.Println("broadcasting termination")
		close(done)
	}()

	// Group goroutine which will write to worklist or
	// start goroutines writing to worklist
	//
	// We can then wait for them all to stop before
	// closing the worklist.
	var wg sync.WaitGroup

	// Create 20 crawler goroutines to fetch each unseen
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			println("started!")
			// TODO: seems the WaitGroup does its job (worklist closed)
			// but perhaps this could be rewritten with a select?
			defer wg.Done()

			for {
				select {
				case link := <-unseenLinks:
					foundLinks := crawl(link)

					wg.Add(1)
					go func() {
						defer wg.Done()

						select {
						case <-done:
							return
						case worklist <- foundLinks:
						}
					}()
				case <-done:
					println("terminating")
					return
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(worklist)
		// Superfluous: if we're here, then the
		// 20 crawlers are all terminated, and the main
		// goroutine (below) would have broke away from
		// sending to unseenLink because of the <-done
		// broadcast
		//close(unseenLinks)
	}()

	// The main goroutine de-duplicates worklist items
	// and sends the unseen ones to the crawlers.
	//
	// We decrement n at each loop, because at every iteration,
	// we fetch something from worklist: if we do, it means
	// a link has been processed by a goroutine.
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		println(n)
		select {
		case <-done:
			// start draining and/or wait for worklist to get
			// closed
			for range worklist {
			}
			// we could return here, but we can also
			// wait next time around: if worklist is closed,
			// then the following read will yield ok=false

		// ok == false never shows up
		case list, ok := <-worklist:
			// worklist is closed
			if !ok {
				return
			}
		loop:
			for _, link := range list {
				if !seen[link.url] {
					seen[link.url] = true
					n++
					select {
					case <-done:
						break loop
					case unseenLinks <- link:
					}
				}
			}
		}
	}
}
