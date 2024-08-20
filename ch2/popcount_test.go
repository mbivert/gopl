package popcount

import (
	"testing"
)

// assumes PopCount() to be correct
func TestAllPopCount(t *testing.T) {
	var xs = []uint64{0b00010011, 98233, 91921884, 191991, 232}

	for _, x := range xs {
		if PopCountLoop(x) != PopCount(x) {
			t.Errorf("PopCountLoop(%d) = %d ≠ %d", x, PopCountLoop(x), PopCount(x))
		}
		if PopCountLast1(x) != PopCount(x) {
			t.Errorf("PopCountLast1(%d) = %d ≠ %d", x, PopCountLast1(x), PopCount(x))
		}
		if PopCountFirst0(x) != PopCount(x) {
			t.Errorf("PopCountFirst0(%d) = %d ≠ %d", x, PopCountFirst0(x), PopCount(x))
		}
	}
}
