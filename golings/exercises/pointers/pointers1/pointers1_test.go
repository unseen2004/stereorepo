package main

import "testing"

func increment(p *int, n int) {
	// I AM NOT DONE
}

func TestIncrement(t *testing.T) {
	x := 5
	increment(&x, 3)
	if x != 8 {
		t.Fatalf("after increment(&x, 3), x = %d, want 8", x)
	}
}
