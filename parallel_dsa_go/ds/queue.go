package ds

import "sync"

type Queue[T any] struct {
	mu   sync.RWMutex
	data []T
}

func (q *Queue[T]) Enqueue(v T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.data = append(q.data, v)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.data) == 0 {
		var zero T
		return zero, false
	}
	res := q.data[0]
	q.data[0] = *new(T)
	q.data = q.data[1:]
	return res, true
}

func (q *Queue[T]) Peek() (T, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	if len(q.data) == 0 {
		var zero T
		return zero, false
	}
	return q.data[0], true
}

func (q *Queue[T]) Len() int {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.data)
}

func (q *Queue[T]) IsEmpty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.data) == 0
}
