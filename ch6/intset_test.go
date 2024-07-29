package main

import (
	"testing"
	"slices"
)

func TestLen(t *testing.T) {
	s := &IntSet{make([]uint64, 0)}

	if s.Len() != 0 {
		t.Errorf("Len({}) != 0")
	}

	s.Add(5)

	if s.Len() != 1 {
		t.Errorf("Len({5}) != 1")
	}

	s.AddAll(99, 1024)

	if s.Len() != 3 {
		t.Errorf("Len({5 99 1024}) != 3")
	}
}

func TestRemove(t *testing.T) {
	s := &IntSet{make([]uint64, 0)}

	s.Remove(0)

	if s.Has(0) {
		t.Errorf("0 not in {}")
	}

	s.Add(5)
	s.Remove(3)

	if s.Has(3) || !s.Has(5) {
		t.Errorf("3 in {5} or 5 not in {5}")
	}

	s.Remove(5)

	if s.Has(5) {
		t.Errorf("5 in {}")
	}

	s.Add(1024)
	if !s.Has(1024) {
		t.Errorf("1024 not in {1024}")
	}

	s.Remove(1024)

	if s.Has(1024) {
		t.Errorf("1024 in {}")
	}
}

func TestClear(t *testing.T) {
	s := &IntSet{make([]uint64, 0)}

	s.AddAll(5, 99, 1024)

	if s.Len() != 3 {
		t.Errorf("Len({5 99 1024}) != 3")
	}

	s.Clear()

	if s.Len() != 0 {
		t.Errorf("Clear() should have deleted everything")
	}
}

func TestCopy(t *testing.T) {
	s := &IntSet{make([]uint64, 0)}

	s.AddAll(5, 99, 1024)

	if s.Len() != 3 {
		t.Errorf("Len({5 99 1024}) != 3")
	}

	x := s.Copy()

	if x.String() != s.String() {
		t.Errorf("Copied set should have the same elements")
	}

	s.Remove(99)

	if s.Has(99) {
		t.Errorf("99 should have been deleted")
	}

	if !x.Has(99) {
		t.Errorf("Deletion in original shouldn't have altered the copy")
	}
}

// NOTE: AddAll is indirectly tested here and there
func TestAddAllElems(t *testing.T) {
	s := &IntSet{make([]uint64, 0)}

	ns := []int{1, 2, 10, 98}
	s.AddAll(ns...)
	ms := s.Elems()

	if !slices.Equal(ns, ms) {
		t.Errorf("%v != %v", ns, ms)
	}
}

func TestIntersectWith(t *testing.T) {
	s := &IntSet{make([]uint64, 0)}
	s.AddAll(1, 42, 18, 67, 910)

	s0 := s.String()
	s.IntersectWith(s)

	if s.String() != s0 {
		t.Errorf("Intersection with self shouldn't change anything")
	}

	x := &IntSet{make([]uint64, 0)}
	x.AddAll(1, 42, 11, 912)

	s.IntersectWith(x)

	if !slices.Equal(s.Elems(), []int{1, 42}) {
		t.Errorf("Intersect {1 42 18 67 910} {1 42 11 912} != {1 42}")
	}
}

func TestDifferenceWith(t *testing.T) {
	s := &IntSet{make([]uint64, 0)}
	s.AddAll(1, 42, 18, 67, 910)

	s.DifferenceWith(s)

	if s.Len() != 0 {
		t.Errorf("{...}-{...} != {}")
	}

	s.AddAll(1, 42, 18, 67, 910)

	x := &IntSet{make([]uint64, 0)}
	x.AddAll(1, 42, 11, 912)

	s.DifferenceWith(x)

	if !slices.Equal(s.Elems(), []int{18, 67, 910}) {
		t.Errorf("{1 42 18 67 910} - {1 42 11 912} != {18 67 910}")
	}

}

func TestSymmetricDifferenceWith(t *testing.T) {
	s := &IntSet{make([]uint64, 0)}
	s.AddAll(1, 42, 18, 67, 910)

	s.SymmetricDifferenceWith(s)

	if s.Len() != 0 {
		t.Errorf("{...} ^- {...} != {}")
	}

	s.AddAll(1, 42, 18, 67, 910)

	x := &IntSet{make([]uint64, 0)}
	x.AddAll(1, 42, 11, 912, 1024)

	s.SymmetricDifferenceWith(x)

	if !slices.Equal(s.Elems(), []int{11, 18, 67, 910, 912, 1024}) {
		t.Errorf("{1 42 18 67 910} ^- {1 42 11 912 1024} != {11 18 67 910 912 1024}")
	}
}
