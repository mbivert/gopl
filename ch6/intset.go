package main

import (
	"bytes"
	"fmt"
)

const (
	// Magic spell.
	intLen = 32 << (^uint(0) >> 63)
)

// i-th bit set <=> i is in the set.
// (words is to be understood (accessed) as a "flat",
// "endless" bit array)
type IntSet struct {
	words []uint
}

func (s *IntSet) Has(x int) bool {
	// word index, bit index
	word, bit := x/intLen, uint(x%intLen)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

func (s *IntSet) Add(x int) {
	word, bit := x/intLen, uint(x%intLen)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= (1<<bit)
}

func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range(t.words) {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

func (s *IntSet) String() string {
	var buf bytes.Buffer

	buf.WriteByte('{')

	for i, word := range(s.words) {
		for j := 0; j < intLen; j++ {
			if word & (1<<j) != 0 {
				// We're beyond the opening {
				if buf.Len() > 1 {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", i*intLen+j)
			}
		}
	}

	buf.WriteByte('}')

	return buf.String()
}

func (s *IntSet) Len() int {
	n := 0

	for _, word := range s.words {
		for j := 0; j < intLen; j++ {
			if word & (1<<j) != 0 {
				n++
			}
		}
	}
	return n
}

func (s *IntSet) Remove(x int) {
	word, bit := x/intLen, x%intLen
	// not here for sure
	if word >= len(s.words) {
		return
	}

	// FTR, "&^" is the so-called "bit clear" operator in Go.
	s.words[word] &^= (1<<bit)
}

func (s *IntSet) Clear() {
	s.words = make([]uint, 0)
}

func (s *IntSet) Copy() *IntSet {
	t := &IntSet{make([]uint, len(s.words))}
	for i := range s.words {
		t.words[i] = s.words[i]
	}
	return t
}

func (s *IntSet) AddAll(ns ...int) {
	for _, n := range ns {
		s.Add(n)
	}
}

func (s *IntSet) Elems() []int {
	ns := make([]int, 0)

	for i, word := range(s.words) {
		for j := 0; j < intLen; j++ {
			if word & (1<<j) != 0 {
				ns = append(ns, i*intLen+j)
			}
		}
	}

	return ns
}

func (s *IntSet) IntersectWith(t *IntSet) {
	for i, tword := range(t.words) {
		if i < len(s.words) {
			s.words[i] &= tword
		}
		// else, only in t, so clearly not in the intersection
	}
}

// Elements which are in s but not in t
func (s *IntSet) DifferenceWith(t *IntSet) {
	for i, tword := range(t.words) {
		if i < len(s.words) {
			s.words[i] &^= tword
		}
		// else, not even in s to begin with
	}
}

// Elements which are either only in s or only in t
func (s *IntSet) SymmetricDifferenceWith(t *IntSet) {
	for i, tword := range(t.words) {
		if i < len(s.words) {
			s.words[i] ^= tword

		// All those are only in t
		} else {
			s.words = append(s.words, tword)
		}
	}
}

func main() {
	s := &IntSet{make([]uint, 0)}

	s.AddAll(5, 19, 42, 1999110232)

	fmt.Println(s.String())
	fmt.Println(intLen)
}
