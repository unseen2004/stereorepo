package main

import "testing"

type rectangle struct {
	width  float64
	height float64
}

func (r rectangle) area() float64 {
	// I AM NOT DONE
	return 0
}

func TestArea(t *testing.T) {
	r := rectangle{width: 3, height: 4}
	if got := r.area(); got != 12 {
		t.Fatalf("area of %+v = %v, want 12", r, got)
	}
}
