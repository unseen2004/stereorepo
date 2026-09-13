package main

import "testing"

func wordCount(words []string) map[string]int {
	// I AM NOT DONE
	return nil
}

func TestWordCount(t *testing.T) {
	got := wordCount([]string{"a", "b", "a"})
	if len(got) != 2 {
		t.Fatalf("wordCount returned %d entries, want 2", len(got))
	}
	if got["a"] != 2 {
		t.Fatalf(`got["a"] = %d, want 2`, got["a"])
	}
	if got["b"] != 1 {
		t.Fatalf(`got["b"] = %d, want 1`, got["b"])
	}
}
