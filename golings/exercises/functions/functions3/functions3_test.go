package main

import "testing"

func sum(nums ...int) (total int) {
	// I AM NOT DONE
	return 0
}

func TestSum(t *testing.T) {
	if got := sum(1, 2, 3); got != 6 {
		t.Fatalf("sum(1, 2, 3) = %d, want 6", got)
	}
	if got := sum(); got != 0 {
		t.Fatalf("sum() = %d, want 0", got)
	}
}
