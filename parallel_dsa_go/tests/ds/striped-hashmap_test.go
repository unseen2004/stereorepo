package tests/ds

import (
	"fmt"
	"sync"
	"testing"
)

func TestStripedHashMap(t *testing.T) {
	m := NewStripedHashMap[string, int](16)

	// Basic Put and Get
	m.Put("a", 1)
	if v, ok := m.Get("a"); !ok || v != 1 {
		t.Errorf("expected 1, got %v (ok: %v)", v, ok)
	}

	// Overwrite
	m.Put("a", 2)
	if v, ok := m.Get("a"); !ok || v != 2 {
		t.Errorf("expected 2, got %v (ok: %v)", v, ok)
	}

	// Get non-existent
	if _, ok := m.Get("b"); ok {
		t.Error("expected ok=false for non-existent key")
	}

	// Delete
	m.Delete("a")
	if _, ok := m.Get("a"); ok {
		t.Error("expected ok=false after delete")
	}

	// Size
	m.Put("x", 10)
	m.Put("y", 20)
	if s := m.Size(); s != 2 {
		t.Errorf("expected size 2, got %d", s)
	}

	// Clear
	m.Clear()
	if s := m.Size(); s != 0 {
		t.Errorf("expected size 0 after clear, got %d", s)
	}
}

func TestStripedHashMapConcurrency(t *testing.T) {
	const (
		numG      = 100
		numOpsPer = 1000
		stripes   = 16
	)
	m := NewStripedHashMap[int, int](stripes)
	var wg sync.WaitGroup
	wg.Add(numG)

	for i := 0; i < numG; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOpsPer; j++ {
				key := (id * numOpsPer) + j
				m.Put(key, key)
				val, ok := m.Get(key)
				if !ok || val != key {
					// Using t.Error inside goroutine is okay but won't stop execution
					fmt.Printf("Error: key %d not found or incorrect value\n", key)
				}
			}
		}(i)
	}

	wg.Wait()
	if s := m.Size(); s != numG*numOpsPer {
		t.Errorf("expected final size %d, got %d", numG*numOpsPer, s)
	}
}
