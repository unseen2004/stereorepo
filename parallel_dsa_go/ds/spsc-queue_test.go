package ds

import (
	"sync"
	"testing"
)

func TestSPSCQueue(t *testing.T) {
	q := NewSPSCQueue[int](4)

	// Single-threaded test
	if !q.Enqueue(1) {
		t.Error("Expected successful enqueue")
	}
	if !q.Enqueue(2) {
		t.Error("Expected successful enqueue")
	}
	if !q.Enqueue(3) {
		t.Error("Expected successful enqueue")
	}
	if !q.Enqueue(4) {
		t.Error("Expected successful enqueue")
	}
	if q.Enqueue(5) {
		t.Error("Expected failed enqueue (full)")
	}

	v, ok := q.Dequeue()
	if !ok || v != 1 {
		t.Errorf("Expected 1, got %v", v)
	}

	if !q.Enqueue(5) {
		t.Error("Expected successful enqueue after dequeue")
	}

	// Multi-threaded SPSC test
	q2 := NewSPSCQueue[int](1024)
	numOps := 1000000
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < numOps; i++ {
			for !q2.Enqueue(i) {
			}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < numOps; i++ {
			var v int
			var ok bool
			for {
				if v, ok = q2.Dequeue(); ok {
					break
				}
			}
			if v != i {
				t.Errorf("Expected %d, got %d", i, v)
			}
		}
	}()

	wg.Wait()
}

func TestSPSCQueueBatch(t *testing.T) {
	q := NewSPSCQueue[int](16)

	// Test EnqueueBatch
	in := []int{1, 2, 3, 4, 5}
	added := q.EnqueueBatch(in)
	if added != 5 {
		t.Errorf("Expected 5 elements added, got %d", added)
	}

	// Test DequeueBatch
	out := make([]int, 3)
	removed := q.DequeueBatch(out)
	if removed != 3 {
		t.Errorf("Expected 3 elements removed, got %d", removed)
	}
	if out[0] != 1 || out[1] != 2 || out[2] != 3 {
		t.Errorf("Unexpected batch content: %v", out)
	}

	removed = q.DequeueBatch(out)
	if removed != 2 {
		t.Errorf("Expected 2 elements removed, got %d", removed)
	}
	if out[0] != 4 || out[1] != 5 {
		t.Errorf("Unexpected batch content: %v", out)
	}

	// Concurrent batch test
	q2 := NewSPSCQueue[int](1024)
	numOps := 1000000
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		batch := make([]int, 64)
		for i := 0; i < numOps; {
			for j := 0; j < 64 && i+j < numOps; j++ {
				batch[j] = i + j
			}
			count := q2.EnqueueBatch(batch)
			i += count
		}
	}()

	go func() {
		defer wg.Done()
		out := make([]int, 64)
		for i := 0; i < numOps; {
			count := q2.DequeueBatch(out)
			for j := 0; j < count; j++ {
				if out[j] != i+j {
					t.Errorf("Expected %d, got %d at index %d", i+j, out[j], j)
				}
			}
			i += count
		}
	}()

	wg.Wait()
}
