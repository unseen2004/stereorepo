package main

import "testing"

func TestZeroValues(t *testing.T) {
	// I AM NOT DONE

	if i != 0 {
		t.Fatalf("zero value of int = %v, want 0", i)
	}
	if f != 0.0 {
		t.Fatalf("zero value of float64 = %v, want 0.0", f)
	}
	if b != false {
		t.Fatalf("zero value of bool = %v, want false", b)
	}
	if s != "" {
		t.Fatalf("zero value of string = %q, want \"\"", s)
	}
}
