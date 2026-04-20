package ds

import (
	"sync"
	"testing"
)

func TestTreiberStack(t *testing.T) {
	s := &TreiberStack[int]{}

	// Single-threaded test
	s.Push(1)
	s.Push(2)
	v, ok := s.Pop()
	if !ok || v != 2 {
		t.Errorf("Expected 2, got %v", v)
	}
	v, ok = s.Pop()
	if !ok || v != 1 {
		t.Errorf("Expected 1, got %v", v)
	}
	_, ok = s.Pop()
	if ok {
		t.Error("Expected empty stack")
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
				s.Push(id*numOps + j)
			}
		}(i)
	}
	wg.Wait()

	count := 0
	for {
		_, ok := s.Pop()
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

func TestTreiberStackConcurrent(t *testing.T) {
	s := &TreiberStack[int]{}
	numGoroutines := 10
	numOps := 10000
	var wg sync.WaitGroup

	wg.Add(numGoroutines * 2)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				s.Push(id*numOps + j)
			}
		}(i)

		go func() {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				for {
					if _, ok := s.Pop(); ok {
						break
					}
				}
			}
		}()
	}

	wg.Wait()
	if _, ok := s.Pop(); ok {
		t.Error("Expected empty stack at the end")
	}
}
