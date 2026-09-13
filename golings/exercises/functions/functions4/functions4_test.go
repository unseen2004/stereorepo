package main

import "testing"

func makeCounter() func() int {
	// I AM NOT DONE
	count := 0
	return func() int {
		return count
	}
}

func TestCounter(t *testing.T) {
	counter := makeCounter()
	if got := counter(); got != 1 {
		t.Fatalf("first call returned %d, want 1", got)
	}
	if got := counter(); got != 2 {
		t.Fatalf("second call returned %d, want 2", got)
	}
	if got := counter(); got != 3 {
		t.Fatalf("third call returned %d, want 3", got)
	}
}
