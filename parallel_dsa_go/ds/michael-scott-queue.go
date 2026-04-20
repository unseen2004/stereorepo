package ds

import (
	"sync/atomic"
)

// Node represents a single element in the queue.
type QueueNode[T any] struct {
	Value T
	Next  atomic.Pointer[QueueNode[T]]
}

// MichaelScottQueue is a lock-free concurrent queue.
type MichaelScottQueue[T any] struct {
	head atomic.Pointer[QueueNode[T]]
	tail atomic.Pointer[QueueNode[T]]
}

// NewMichaelScottQueue creates a new queue with a dummy node.
func NewMichaelScottQueue[T any]() *MichaelScottQueue[T] {
	dummy := &QueueNode[T]{}
	q := &MichaelScottQueue[T]{}
	q.head.Store(dummy)
	q.tail.Store(dummy)
	return q
}

// Enqueue adds a value to the end of the queue.
func (q *MichaelScottQueue[T]) Enqueue(v T) {
	n := &QueueNode[T]{Value: v}
	for {
		tail := q.tail.Load()
		next := tail.Next.Load()

		// Check if tail is still consistent
		if tail == q.tail.Load() {
			if next == nil {
				// Try to link the new node to the end of the list
				if tail.Next.CompareAndSwap(nil, n) {
					// Try to advance the tail pointer (even if it fails, others will help)
					q.tail.CompareAndSwap(tail, n)
					return
				}
			} else {
				// Tail is lagging, try to advance it for the other thread
				q.tail.CompareAndSwap(tail, next)
			}
		}
	}
}

// Dequeue removes and returns a value from the front of the queue.
func (q *MichaelScottQueue[T]) Dequeue() (T, bool) {
	for {
		head := q.head.Load()
		tail := q.tail.Load()
		next := head.Next.Load()

		// Check if head is still consistent
		if head == q.head.Load() {
			if head == tail {
				if next == nil {
					// Queue is empty
					var zero T
					return zero, false
				}
				// Tail is lagging, try to advance it
				q.tail.CompareAndSwap(tail, next)
			} else {
				// Read value before CAS, as the node might be reused/modified after CAS
				v := next.Value
				if q.head.CompareAndSwap(head, next) {
					return v, true
				}
			}
		}
	}
}
