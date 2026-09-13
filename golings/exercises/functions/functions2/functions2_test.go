package main

import "testing"

func divmod(a, b int) (int, int) {
	// I AM NOT DONE
	return 0, 0
}

func TestDivmod(t *testing.T) {
	q, r := divmod(7, 3)
	if q != 2 || r != 1 {
		t.Fatalf("divmod(7, 3) = %d, %d, want 2, 1", q, r)
	}
}
