package main

import (
	"testing"
	"time"
)

func parallelSum(nums []int) int {
	// I AM NOT DONE
	return 0
}

func TestParallelSum(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	done := make(chan int)
	go func() {
		done <- parallelSum(nums)
	}()
	select {
	case got := <-done:
		if got != 15 {
			t.Fatalf("parallelSum(%v) = %d, want 15", nums, got)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out - parallelSum never returned")
	}
}
