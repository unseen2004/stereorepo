package main

import "testing"

func grade(score int) string {
	// I AM NOT DONE
	return ""
}

func TestGrade(t *testing.T) {
	cases := map[int]string{
		95: "A",
		85: "B",
		75: "C",
		65: "F",
	}
	for score, want := range cases {
		if got := grade(score); got != want {
			t.Fatalf("grade(%d) = %q, want %q", score, got, want)
		}
	}
}
