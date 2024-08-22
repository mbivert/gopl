package main

/*
 * At least up until there, linear in the payload size;
 * see ./run-pingpong.sh

10             : 11939226         byte/sec
100            : 123211300        byte/sec
1000           : 1223219666       byte/sec
10000          : 12510513333      byte/sec
100000         : 121146966666     byte/sec
1000000        : 1223674333333    byte/sec
10000000       : 12493956666666   byte/sec

 */

import (
	"flag"
	"fmt"
	"time"
)

var (
	sFlag = flag.Int("s", 15, "payload size (bytes)")
	rFlag = flag.Int("r", 5, "number of run (0 to disable)")
)

func mkMsg(n int) []byte {
	msg := make([]byte, n)
	for i := 0; i < n; i++ {
		msg[i] = 'a'
	}
	return msg
}

func main() {
	flag.Parse()

	a := make(chan []byte)
	b := make(chan []byte)
	msg := mkMsg(*sFlag)

	start := time.Now()

	go func() {
		for m := range a {
			b <- m
		}
	}()

	a <- msg

	var n, r int
	for m := range b {
		// 1sec should absorb the instrumentation cost enough.
		if time.Now().After(start.Add(1 * time.Second)) {
			fmt.Printf("%d byte/sec\n", n*len(msg))
			start = time.Now()
			n = 0
			r++
			if *rFlag != 0 && r == *rFlag {
				break
			}
		}
		n++
		a <- m
	}
	close(a)
	close(b)
}
