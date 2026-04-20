package ds

import "sync"

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

type Counter[T Number] struct {
	mu sync.RWMutex
	v  T
}

func (c *Counter[T]) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v++
}

func (c *Counter[T]) Dec() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v--
}

func (c *Counter[T]) Get() T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.v
}

func (c *Counter[T]) Update(v T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v = v
}

func (c *Counter[T]) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.v = 0
}
