package main

import (
	"os"
	"sync"
)

var loadPcOnce sync.Once

// pc[i] is the population count of i.
// i.e. number of bits set to 1 in the "integer" i
var pc [256]byte

func loadPc() {
	for i := range pc {
		// i/2 <=> i >> 1;
		//	-> pc[i/2] is the population count of i but its last bit
		// i&1 <=> i's last bit = 1
		// -> This is a "recursive" initialization, which relies
		// on the fact that pc[0] = 0 (default integer/byte init value)
		pc[i] = pc[i/2] + byte(i&1)
	}
}

// PopCount returns the population count (number of set bits) of x.
func PopCount(x uint64) int {
	loadPcOnce.Do(loadPc)

	// split x in 256 bytes chunks, and use pc[] to compute
	// the number of bits set in each chunks; sum
	return int(pc[byte(x>>(0*8))] +
		pc[byte(x>>(1*8))] +
		pc[byte(x>>(2*8))] +
		pc[byte(x>>(3*8))] +
		pc[byte(x>>(4*8))] +
		pc[byte(x>>(5*8))] +
		pc[byte(x>>(6*8))] +
		pc[byte(x>>(7*8))])
}

func main() {
	// quick test
	if PopCount(0b010010101) == 4 {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
