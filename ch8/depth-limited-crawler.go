package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

var depth int
var ncrawl int

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

func extractLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
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
	go func() { worklist <- mkLinks(0, flag.Args()) }()

	// Create 20 crawler goroutines to fetch each unseen
	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link)
				go func() { worklist <- foundLinks }()
			}
		}()
	}

	// The main goroutine de-duplicates worklist items
	// and sends the unseen ones to the crawlers.
	//
	// We decrement n at each loop, because at every iteration,
	// we fetch something from worklist: if we do, it means
	// a link has been processed by a goroutine.
	seen := make(map[string]bool)
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			if !seen[link.url] {
				seen[link.url] = true
				n++
				unseenLinks <- link
			}
		}
	}
}
