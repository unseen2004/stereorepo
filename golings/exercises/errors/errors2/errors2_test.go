package main

import (
	"fmt"
	"math"
	"testing"
)

func sqrt(x float64) (float64, error) {
	// I AM NOT DONE
	if x < 0 {
		return 0, nil
	}
	return math.Sqrt(x), nil
}

func TestSqrt(t *testing.T) {
	if _, err := sqrt(-1); err == nil {
		t.Fatal("sqrt(-1) returned no error, want an error")
	}
	got, err := sqrt(9)
	if err != nil {
		t.Fatalf("sqrt(9) returned error %v", err)
	}
	if got != 3 {
		t.Fatalf("sqrt(9) = %v, want 3", got)
	}
}
