package main

/*
 * This is just a prototype: at least we don't:
 *	- fetch non-HTML files, including images, fonts, PDFs, etc.;
 *	(this is probably just a matter of checking the response's Content-Type)
 *	- adjust image's URLs;
 *
 * Usage is e.g.:
 *	go run mirror.go -scheme="file" -host="" -path="$PWD/tales" https://tales.mbivert.com/ tales/
 *
 * This will changes <a>'s href to point to file:/// URLs: note that we're not
 * directly linking to index.html, which we perhaps should do in case
 * the scheme is file:// (perhaps this could even be an option).
 *
 * NOTE: for some reason, html.Parse() thinks it can parse images? e.g.
 *	process("https://tales.mbivert.com/bargue-plate-I-30.jpg", os.Stdout)
 */

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"golang.org/x/net/html"
)

var ncrawl int          // number of concurrent crawler
var mainHost string     // main hostname, computed from the CLI provided URL
var mirrorScheme string // mirror scheme (e.g. file, https, etc.)
var mirrorPath string   // mirror "root" path
var mirrorHost string   // mirror hostname, used to rebuild the "internal" URLs
var dir string          // output directory

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

func isIntURL(url *url.URL) bool {
	return url.Hostname() == mainHost
}

func process(url string, fh io.Writer) ([]string, error) {
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

	fmt.Println("Parsed")

	var links []string
	visitNode := func(n *html.Node) {
		switch n.Type {
		case html.TextNode:
			fh.Write([]byte(n.Data))
		case html.DocumentNode:
			// nuthin'
		case html.ElementNode:
			fmt.Fprintf(fh, "<%s", n.Data)

			for _, a := range n.Attr {
				if n.Data == "a" && a.Key == "href" {
					link, err := resp.Request.URL.Parse(a.Val)
					if err != nil {
						continue // ignore bad URLs
					}
					if isIntURL(link) {
						// Consider URLs with different fragments to be identical;
						// perhaps it'd make sense to ignore URLs containing a
						// query as well.
						link.Fragment = ""

						links = append(links, link.String())
					}
					link.Scheme = mirrorScheme
					link.Path = filepath.Join(mirrorPath, link.Path)
					link.Host = mirrorHost
					a.Val = link.String()
				}
				fmt.Fprintf(fh, " %s=\"%s\"", a.Key, a.Val)
			}
			fmt.Fprintf(fh, ">")
		// untested...
		case html.CommentNode:
			fmt.Fprintf(fh, "<-- %s -->", n.Data)
		case html.DoctypeNode:
			fh.Write([]byte(n.Data))
		}
	}
	closeTags := func(n *html.Node) {
		switch n.Type {
		case html.TextNode:
		case html.DocumentNode:
			// nuthin'
		case html.ElementNode:
			fmt.Fprintf(fh, "</%s>", n.Data)
		case html.CommentNode:
		case html.DoctypeNode:
		}
	}

	forEachNode(doc, visitNode, closeTags)
	return links, nil
}

// copmute (local) path where to save the (altered) content
// of the given URL.
func getPath(link string) (string, error) {
	// NOTE: we're parsing each URL twice then
	url, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, url.Path)
	if url.Path[len(url.Path)-1] == '/' {
		path = filepath.Join(path, "index.html")
	}
	return path + url.RawQuery, nil
}

func openFile(path string) (*os.File, error) {
	// make sure intermediate directories all exist
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	// try to create the file
	return os.Create(path)
}

func crawl(link string) []string {
	path, err := getPath(link)
	if err != nil {
		log.Print(err)
		return []string{}
	}

	fh, err := openFile(path)
	if err != nil {
		log.Print(err)
		return []string{}
	}
	defer fh.Close()

	fmt.Printf("%-70s -> %s\n", link, path)

	links, err := process(link, fh)
	if err != nil {
		log.Print(err)
	}
	return links
}

// note: things stale if there's an error in crawl() on the first URL (at least)
func main() {
	flag.IntVar(&ncrawl, "ncrawl", 20, "number of concurrent crawlers")
	flag.StringVar(&mirrorScheme, "scheme", "https://", "Scheme for mirror")
	flag.StringVar(&mirrorHost, "host", "", "mirror hostname")
	flag.StringVar(&mirrorPath, "path", "", "mirror path")

	flag.Parse()

	if len(flag.Args()) != 2 {
		log.Fatal("%s [-ncrawl=n] <url> <path/to/directory>", os.Args[0])
	}

	var link string
	link, dir = flag.Args()[0], flag.Args()[1]

	url, err := url.Parse(link)
	if err != nil {
		log.Fatal(err)
	}

	mainHost = url.Hostname()

	//	process("https://tales.mbivert.com/bargue-plate-I-30.jpg", os.Stdout)
	//	os.Exit(0)

	var n int                        // number of pending links to unseenLinks
	worklist := make(chan []string)  // lists of URLs, may have duplicates
	unseenLinks := make(chan string) // de-duplicated URLs

	// Add command-line arguments to worklist.
	n = len(flag.Args())
	go func() { worklist <- []string{link} }()

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
			if !seen[link] {
				seen[link] = true
				n++
				unseenLinks <- link
			}
		}
	}
}
