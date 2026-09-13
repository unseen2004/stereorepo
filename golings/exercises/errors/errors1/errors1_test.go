package main

import (
	"errors"
	"testing"
)

func divide(a, b int) (int, error) {
	// I AM NOT DONE
	return a / b, nil
}

func TestDivide(t *testing.T) {
	if _, err := divide(1, 0); err == nil {
		t.Fatal("divide(1, 0) returned no error, want an error")
	}
	q, err := divide(6, 2)
	if err != nil {
		t.Fatalf("divide(6, 2) returned error %v", err)
	}
	if q != 3 {
		t.Fatalf("divide(6, 2) = %d, want 3", q)
	}
}
