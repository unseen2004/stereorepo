package ds

import (
	"sync/atomic"
)

type Node[T any] struct {
	Value T
	Next  *Node[T]
}

type TreiberStack[T any] struct {
	head atomic.Pointer[Node[T]]
}

func (s *TreiberStack[T]) Push(v T) {
	n := &Node[T]{Value: v}
	for {
		old := s.head.Load()
		n.Next = old
		if s.head.CompareAndSwap(old, n) {
			return
		}
	}
}

func (s *TreiberStack[T]) Pop() (T, bool) {
	for {
		old := s.head.Load()
		if old == nil {
			var zero T
			return zero, false
		}
		next := old.Next
		if s.head.CompareAndSwap(old, next) {
			return old.Value, true
		}
	}
}
