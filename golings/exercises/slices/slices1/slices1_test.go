package main

import "testing"

func buildList(n int) []int {
	// I AM NOT DONE
	return nil
}

func TestBuildList(t *testing.T) {
	got := buildList(3)
	if len(got) != 3 {
		t.Fatalf("buildList(3) has length %d, want 3", len(got))
	}
	for i, want := range []int{1, 2, 3} {
		if got[i] != want {
			t.Fatalf("buildList(3)[%d] = %d, want %d", i, got[i], want)
		}
	}
}
