package popcount

// pc[i] is the population count of i.
// i.e. number of bits set to 1 in the "integer" i
var pc [256]byte

func init() {
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
	// split x in 256 bits chunks, and use pc[] to compute
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

func PopCountLoop(x uint64) int {
	var r int

	for i := 0; i < 8; i++ {
		r += int(pc[byte(x>>(i*8))])
	}

	return r
}

func PopCountLast1(x uint64) int {
	var r int

	for i := 0; i < 64; i++ {
		if (x>>i)&1 == 1 {
			r++
		}
	}

	return r
}

func PopCountFirst0(x uint64) int {
	var r int

	for x != 0 {
		r++
		// « x&(x-1) clears the rightmost non-zero bit of x »
		x = x & (x - 1)
	}

	return r

}
