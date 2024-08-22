package main

import (
	"flag"
	"fmt"
	"time"
)

/*
 * See ./run-maxgo.sh

Time to go through the pipeline:
	main received 42         with 1               goroutines, in 5.231µs
	main received 42         with 10              goroutines, in 30.445µs
	main received 42         with 100             goroutines, in 156.982µs
	main received 42         with 1000            goroutines, in 1.388409ms
	main received 42         with 10000           goroutines, in 12.624567ms
	main received 42         with 100000          goroutines, in 115.815265ms
	main received 42         with 1000000         goroutines, in 1.233440883s
	main received 42         with 10000000        goroutines, in 12.600425441s
Creating a huge pipeline:
	main received -1         with 1               goroutines, in 6.243µs
	main received -10        with 10              goroutines, in 19.876µs
	main received -100       with 100             goroutines, in 194.183µs
	main received -1000      with 1000            goroutines, in 1.697617ms
	main received -10000     with 10000           goroutines, in 21.696693ms
	main received -100000    with 100000          goroutines, in 271.474083ms
	main received -1000000   with 1000000         goroutines, in 2.680637863s

*/

var n int

var cancel = make(chan struct{})

var (
	mFlag = flag.Int("m", 0, "If non-zero, build a m-stage pipeline")
	vFlag = flag.Bool("v", false, "Verbose output")
	// Meaning, create a pipeline with *mFlag goroutines.
	nFlag = flag.Bool("n", false, "Don't send a value through the pipeline")
)

func launchGo(in, endp chan int) {
	var out chan int
	var v int

	n++
	if *vFlag {
		fmt.Printf("%d\n", n)
	}

	if *mFlag != 0 && n < *mFlag {
		out = make(chan int)
		go launchGo(out, endp)

		// Last stage of the pipeline through which
		// no value is going through: notify main that
		// we can cancel everyone.
	} else if *nFlag {
		endp <- -n
		return
	}

	// Otherwise, wait either to be cancelled or to
	// receive a value
	select {
	case v = <-in:
		if out != nil {
			out <- v
		}
	case <-cancel:
	}

	if *vFlag {
		fmt.Printf("%-10d received %v\n", n, v)
	}
	if out == nil {
		endp <- v
	}
}

func main() {
	flag.Parse()

	c := make(chan int)
	start := time.Now()
	endp := make(chan int)

	go launchGo(c, endp)

	if !*nFlag {
		c <- 42
	}

	select {
	case v := <-endp:
		fmt.Printf("main received %-10d with %-15d goroutines, in %s\n", v, *mFlag, time.Since(start))
		close(cancel)
	}
}
