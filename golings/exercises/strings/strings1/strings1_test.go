package main

import "testing"

func reverse(s string) string {
	// I AM NOT DONE
	return ""
}

func TestReverse(t *testing.T) {
	if got := reverse("hello"); got != "olleh" {
		t.Fatalf(`reverse("hello") = %q, want "olleh"`, got)
	}
	if got := reverse("héllo"); got != "olléh" {
		t.Fatalf(`reverse("héllo") = %q, want "olléh"`, got)
	}
}
