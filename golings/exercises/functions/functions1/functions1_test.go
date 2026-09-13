package main

import "testing"

func add(a, b int) int {
	// I AM NOT DONE
	return 0
}

func TestAdd(t *testing.T) {
	if got := add(2, 3); got != 5 {
		t.Fatalf("add(2, 3) = %d, want 5", got)
	}
}
