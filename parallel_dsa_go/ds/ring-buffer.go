package ds

import "sync"

type RingBuffer[T any] struct {
	mu       sync.RWMutex
	data     []T
	head     int
	tail     int
	count    int
	capacity int
}

func NewRingBuffer[T any](capacity int) *RingBuffer[T] {
	if capacity <= 0 {
		capacity = 16
	}
	return &RingBuffer[T]{
		data:     make([]T, capacity),
		capacity: capacity,
	}
}

func (r *RingBuffer[T]) Enqueue(v T) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.count == r.capacity {
		return false
	}

	r.data[r.tail] = v
	r.tail = (r.tail + 1) % r.capacity
	r.count++
	return true
}

func (r *RingBuffer[T]) Dequeue() (T, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.count == 0 {
		var zero T
		return zero, false
	}

	v := r.data[r.head]
	r.data[r.head] = *new(T)
	r.head = (r.head + 1) % r.capacity
	r.count--
	return v, true
}

func (r *RingBuffer[T]) Peek() (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.count == 0 {
		var zero T
		return zero, false
	}

	return r.data[r.head], true
}

func (r *RingBuffer[T]) IsFull() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.count == r.capacity
}

func (r *RingBuffer[T]) IsEmpty() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.count == 0
}

func (r *RingBuffer[T]) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.count
}

func (r *RingBuffer[T]) Cap() int {
	return r.capacity
}

func (r *RingBuffer[T]) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data = make([]T, r.capacity)
	r.head = 0
	r.tail = 0
	r.count = 0
}
