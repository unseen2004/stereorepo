package main

import "testing"

func sumTo(n int) int {
	// I AM NOT DONE
	return 0
}

func TestSumTo(t *testing.T) {
	if got := sumTo(5); got != 15 {
		t.Fatalf("sumTo(5) = %d, want 15", got)
	}
	if got := sumTo(0); got != 0 {
		t.Fatalf("sumTo(0) = %d, want 0", got)
	}
}
