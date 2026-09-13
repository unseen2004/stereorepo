package main

import (
	"strings"
	"testing"
)

func countVowels(s string) int {
	// I AM NOT DONE
	return 0
}

func TestCountVowels(t *testing.T) {
	if got := countVowels("Hello World"); got != 3 {
		t.Fatalf(`countVowels("Hello World") = %d, want 3`, got)
	}
	if got := countVowels("sky"); got != 0 {
		t.Fatalf(`countVowels("sky") = %d, want 0`, got)
	}
}
