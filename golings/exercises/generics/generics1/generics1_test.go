package main

import "testing"

func Max[T int | float64](a, b T) T {
	// I AM NOT DONE
	var zero T
	return zero
}

func TestMax(t *testing.T) {
	if got := Max(3, 7); got != 7 {
		t.Fatalf("Max(3, 7) = %d, want 7", got)
	}
	if got := Max(2.5, 1.5); got != 2.5 {
		t.Fatalf("Max(2.5, 1.5) = %v, want 2.5", got)
	}
}
