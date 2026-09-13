package main

import "testing"

func copySlice(src []int) []int {
	// I AM NOT DONE
	return nil
}

func TestCopySlice(t *testing.T) {
	src := []int{1, 2, 3}
	dst := copySlice(src)
	dst[0] = 99
	if src[0] != 1 {
		t.Fatal("dst shares memory with src - the copy is not independent")
	}
	if len(dst) != 3 || dst[1] != 2 || dst[2] != 3 {
		t.Fatalf("copySlice([1 2 3]) = %v, want [99 2 3]", dst)
	}
}
