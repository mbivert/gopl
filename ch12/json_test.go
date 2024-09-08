package main

import (
	"encoding/json"
	"testing"
)

// Probably incomplete, but good enough to deal with strangelove.

func TestIntArray(t *testing.T) {
	a := []int{1, 2, 3}
	xs, errx := json.Marshal(a)
	if errx != nil {
		t.Errorf("Unexpected encoding/json error: %s", errx)
	}
	ys, erry := Marshal(a)
	if erry != nil {
		t.Errorf("Unexpected json error: %s", erry)
	}
	if string(xs) != string(ys) {
		t.Errorf("%s != %s", xs, ys)
	}
}

func TestFloatArray(t *testing.T) {
	a := []float64{1.2, 2.3, 3.4}
	xs, errx := json.Marshal(a)
	if errx != nil {
		t.Errorf("Unexpected encoding/json error: %s", errx)
	}
	ys, erry := Marshal(a)
	if erry != nil {
		t.Errorf("Unexpected json error: %s", erry)
	}
	if string(xs) != string(ys) {
		t.Errorf("%s != %s", xs, ys)
	}
}

func TestBoolArray(t *testing.T) {
	a := []bool{true, false, true}
	xs, errx := json.Marshal(a)
	if errx != nil {
		t.Errorf("Unexpected encoding/json error: %s", errx)
	}
	ys, erry := Marshal(a)
	if erry != nil {
		t.Errorf("Unexpected json error: %s", erry)
	}
	if string(xs) != string(ys) {
		t.Errorf("%s != %s", xs, ys)
	}
}

func TestSruct(t *testing.T) {
	s := struct {
		Foo string
		Bar string
		Num int
	}{"foo", "bar", 7}

	xs, errx := json.Marshal(s)
	if errx != nil {
		t.Errorf("Unexpected encoding/json error: %s", errx)
	}
	ys, erry := Marshal(s)
	if erry != nil {
		t.Errorf("Unexpected json error: %s", erry)
	}
	if string(xs) != string(ys) {
		t.Errorf("%s != %s", xs, ys)
	}
}

func TestMap(t *testing.T) {
	s := map[string]string {
		"foo" : "bare",
		"bar" : "baze",
		"baz" : "fooe",
	}

	xs, errx := json.Marshal(s)
	if errx != nil {
		t.Errorf("Unexpected encoding/json error: %s", errx)
	}
	ys, erry := Marshal(s)
	if erry != nil {
		t.Errorf("Unexpected json error: %s", erry)
	}
	if string(xs) != string(ys) {
		t.Errorf("%s != %s", xs, ys)
	}
}

func TestBigStruct(t *testing.T) {
	xs, errx := json.Marshal(strangelove)
	if errx != nil {
		t.Errorf("Unexpected encoding/json error: %s", errx)
	}
	ys, erry := Marshal(strangelove)
	if erry != nil {
		t.Errorf("Unexpected json error: %s", erry)
	}
	if string(xs) != string(ys) {
		t.Errorf("%s != %s", xs, ys)
	}
}
