package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var verbose = flag.Bool("v", false, "Show progress")

type fSize struct {
	root string
	size int64
}

// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on fileSizes.
func walkDir(root, dir string, wg *sync.WaitGroup, fileSizes chan<- fSize) {
	defer wg.Done()

	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			wg.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(root, subdir, wg, fileSizes)
		} else {
			fi, err := entry.Info()
			if err != nil {
				fmt.Fprintf(os.Stderr, "du: %v\n", err)
			} else {
				fileSizes <- fSize{root, fi.Size()}
			}
		}
	}
}

// 20 simultaneous calls to dirents max
var sema = make(chan struct{}, 20)

// dirents returns the entries of directory dir.
func dirents(dir string) []fs.DirEntry {
	// acquire lock
	sema <- struct{}{}

	// release
	defer func() { <-sema }()

	// NOTE: gopl relies on ioutil.ReadDir(), which is now
	// deprecated.
	xs, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		// xs may not be nil (partial results)
	}
	return xs
}

type Total struct {
	nfiles int64
	nbytes int64
}

func printDiskUsage(totals map[string]*Total) {
	for k, v := range totals {
		fmt.Printf("%30s: %10d files %5.1f GB\n", k, v.nfiles, float64(v.nbytes)/1e9)
	}
}

func main() {
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Traverse the file tree.
	fileSizes := make(chan fSize)
	var wg sync.WaitGroup

	for _, root := range roots {
		wg.Add(1)
		go walkDir(root, root, &wg, fileSizes)
	}

	go func() {
		wg.Wait()
		close(fileSizes)
	}()

	var tick <-chan time.Time
	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}

	// root => {nfiles, nbytes}
	totals := make(map[string]*Total, len(roots))
	for _, root := range roots {
		totals[root] = &Total{0, 0}
	}

loop:
	for {
		select {
		// if -v is not set, then tick is nil, and
		// we'll never read from it.
		case <-tick:
			printDiskUsage(totals)
		case sz, ok := <-fileSizes:
			if !ok {
				break loop
			}
			total := totals[sz.root]
			total.nfiles++
			total.nbytes += sz.size
		}
	}

	printDiskUsage(totals)
}
