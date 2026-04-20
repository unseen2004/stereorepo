package ds

import (
	"sync"
	"testing"
)

func TestMichaelScottQueue(t *testing.T) {
	q := NewMichaelScottQueue[int]()

	// Single-threaded test
	q.Enqueue(1)
	q.Enqueue(2)
	v, ok := q.Dequeue()
	if !ok || v != 1 {
		t.Errorf("Expected 1, got %v", v)
	}
	v, ok = q.Dequeue()
	if !ok || v != 2 {
		t.Errorf("Expected 2, got %v", v)
	}
	_, ok = q.Dequeue()
	if ok {
		t.Error("Expected empty queue")
	}

	// Multi-threaded test
	var wg sync.WaitGroup
	numGoroutines := 100
	numOps := 1000

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				q.Enqueue(id*numOps + j)
			}
		}(i)
	}
	wg.Wait()

	count := 0
	for {
		_, ok := q.Dequeue()
		if !ok {
			break
		}
		count++
	}

	expected := numGoroutines * numOps
	if count != expected {
		t.Errorf("Expected %v elements, got %v", expected, count)
	}
}

func TestMichaelScottQueueConcurrent(t *testing.T) {
	q := NewMichaelScottQueue[int]()
	numGoroutines := 10
	numOps := 10000
	var wg sync.WaitGroup

	wg.Add(numGoroutines * 2)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				q.Enqueue(id*numOps + j)
			}
		}(i)

		go func() {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				for {
					if _, ok := q.Dequeue(); ok {
						break
					}
				}
			}
		}()
	}

	wg.Wait()
	if _, ok := q.Dequeue(); ok {
		t.Error("Expected empty queue at the end")
	}
}
