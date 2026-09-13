package main

import "testing"

func intToFloat(n int) float64 {
	// I AM NOT DONE
	return 0
}

func floatToInt(f float64) int {
	// I AM NOT DONE
	return 0
}

func TestConversions(t *testing.T) {
	if got := intToFloat(5); got != 5.0 {
		t.Fatalf("intToFloat(5) = %v, want 5.0", got)
	}
	if got := floatToInt(3.9); got != 3 {
		t.Fatalf("floatToInt(3.9) = %v, want 3", got)
	}
}
