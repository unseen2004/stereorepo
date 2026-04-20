package ds

import "sync"

type SafeSlice[T any] struct {
	mu   sync.RWMutex
	data []T
}

func (s *SafeSlice[T]) Add(v T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = append(s.data, v)
}

func (s *SafeSlice[T]) TryAdd(v T) bool {
	if s.mu.TryLock() {
		defer s.mu.Unlock()
		s.data = append(s.data, v)
		return true
	}
	return false
}

func (s *SafeSlice[T]) At(i int) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if i < 0 || i >= len(s.data) {
		var zero T
		return zero, false
	}
	return s.data[i], true
}

func (s *SafeSlice[T]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

func (s *SafeSlice[T]) GetAll() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]T, len(s.data))
	copy(res, s.data)
	return res
}

func (s *SafeSlice[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = nil
}
