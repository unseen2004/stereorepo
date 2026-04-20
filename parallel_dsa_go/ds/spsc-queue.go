package ds

import (
	"sync/atomic"
)

type SPSCQueue[T any] struct {
	_ [64]byte

	tail       atomic.Uint64
	cachedHead uint64
	
	_ [64]byte

	head       atomic.Uint64
	cachedTail uint64
	
	_ [64]byte

	data     []T
	capacity uint64
	mask     uint64
	
	_ [64]byte
}

func NewSPSCQueue[T any](size uint64) *SPSCQueue[T] {
	capacity := uint64(1)
	for capacity < size {
		capacity <<= 1
	}

	return &SPSCQueue[T]{
		data:     make([]T, capacity),
		capacity: capacity,
		mask:     capacity - 1,
	}
}

func (q *SPSCQueue[T]) Enqueue(v T) bool {
	tail := q.tail.Load()

	if (tail + 1 - q.cachedHead) > q.capacity {
		q.cachedHead = q.head.Load()
		if (tail + 1 - q.cachedHead) > q.capacity {
			return false
		}
	}

	q.data[tail&q.mask] = v
	q.tail.Store(tail + 1)
	return true
}

func (q *SPSCQueue[T]) EnqueueBatch(in []T) int {
	tail := q.tail.Load()
	batchSize := uint64(len(in))

	if (tail + batchSize - q.cachedHead) > q.capacity {
		q.cachedHead = q.head.Load()
		if (tail + batchSize - q.cachedHead) > q.capacity {
			available := q.capacity - (tail - q.cachedHead)
			if available <= 0 {
				return 0
			}
			batchSize = available
		}
	}

	for i := uint64(0); i < batchSize; i++ {
		q.data[(tail+i)&q.mask] = in[i]
	}

	q.tail.Store(tail + batchSize)
	return int(batchSize)
}

func (q *SPSCQueue[T]) Dequeue() (T, bool) {
	head := q.head.Load()

	if head == q.cachedTail {
		q.cachedTail = q.tail.Load()
		if head == q.cachedTail {
			var zero T
			return zero, false
		}
	}

	v := q.data[head&q.mask]
	var zero T
	q.data[head&q.mask] = zero
	
	q.head.Store(head + 1)
	return v, true
}

func (q *SPSCQueue[T]) DequeueBatch(out []T) int {
	head := q.head.Load()

	if head == q.cachedTail {
		q.cachedTail = q.tail.Load()
		if head == q.cachedTail {
			return 0
		}
	}

	available := q.cachedTail - head
	batchSize := uint64(len(out))
	if available < batchSize {
		batchSize = available
	}

	var zero T
	for i := uint64(0); i < batchSize; i++ {
		idx := (head + i) & q.mask
		out[i] = q.data[idx]
		q.data[idx] = zero
	}

	q.head.Store(head + batchSize)
	return int(batchSize)
}

func (q *SPSCQueue[T]) IsFull() bool {
	tail := q.tail.Load()
	if (tail + 1 - q.cachedHead) > q.capacity {
		q.cachedHead = q.head.Load()
	}
	return (tail + 1 - q.cachedHead) > q.capacity
}

func (q *SPSCQueue[T]) IsEmpty() bool {
	head := q.head.Load()
	if head == q.cachedTail {
		q.cachedTail = q.tail.Load()
	}
	return head == q.cachedTail
}
