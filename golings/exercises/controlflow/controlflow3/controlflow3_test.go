package main

import "testing"

func total(nums []int) int {
	// I AM NOT DONE
	return 0
}

func TestTotal(t *testing.T) {
	nums := []int{4, 5, 6}
	if got := total(nums); got != 15 {
		t.Fatalf("total(%v) = %d, want 15", nums, got)
	}
	if got := total(nil); got != 0 {
		t.Fatalf("total(nil) = %d, want 0", got)
	}
}
