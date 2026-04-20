package ds

import (
	"fmt"
	"hash/fnv"
	"sync"
)

type stripe[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

type StripedHashMap[K comparable, V any] struct {
	stripes []*stripe[K, V]
	count   int
}

func NewStripedHashMap[K comparable, V any](numStripes int) *StripedHashMap[K, V] {
	if numStripes <= 0 {
		numStripes = 16
	}
	shm := &StripedHashMap[K, V]{
		stripes: make([]*stripe[K, V], numStripes),
		count:   numStripes,
	}
	for i := 0; i < numStripes; i++ {
		shm.stripes[i] = &stripe[K, V]{
			data: make(map[K]V),
		}
	}
	return shm
}

func (m *StripedHashMap[K, V]) hash(key K) uint32 {
	h := fnv.New32a()
	// Efficiently handle common types, fallback to fmt.Fprintf for others
	switch v := any(key).(type) {
	case string:
		h.Write([]byte(v))
	default:
		fmt.Fprintf(h, "%v", v)
	}
	return h.Sum32()
}

func (m *StripedHashMap[K, V]) getStripe(key K) *stripe[K, V] {
	h := m.hash(key)
	return m.stripes[h%uint32(m.count)]
}

func (m *StripedHashMap[K, V]) Put(key K, val V) {
	s := m.getStripe(key)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = val
}

func (m *StripedHashMap[K, V]) Get(key K) (V, bool) {
	s := m.getStripe(key)
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

func (m *StripedHashMap[K, V]) Delete(key K) {
	s := m.getStripe(key)
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

func (m *StripedHashMap[K, V]) Size() int {
	total := 0
	for _, s := range m.stripes {
		s.mu.RLock()
		total += len(s.data)
		s.mu.RUnlock()
	}
	return total
}

func (m *StripedHashMap[K, V]) Clear() {
	for _, s := range m.stripes {
		s.mu.Lock()
		s.data = make(map[K]V)
		s.mu.Unlock()
	}
}
