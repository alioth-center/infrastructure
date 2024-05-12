package concurrency

import "sync"

// Slice is a thread-safe slice
type Slice[T any] interface {
	// Append appends item to the end of slice
	// example:
	//
	//	slice := NewSlice[int]()
	//	slice.Append(1)
	//	slice.Append(2)
	// then
	// slice.Items() -> [1, 2]
	Append(item T)

	// Appends appends items to the end of slice
	// example:
	//
	//	slice := NewSlice[int]()
	//	slice.Appends(1, 2)
	// then
	// slice.Items() -> [1, 2]
	Appends(items ...T)

	// SubSlice returns a new slice from start to end, means [start, end)
	// example:
	//
	//	slice := NewSlice[int]()
	//	slice.Appends(1, 2, 3, 4, 5)
	//	slice = slice.SubSlice(1, 3)
	// then
	// slice.Items() -> [2, 3]
	SubSlice(start, end int) Slice[T]

	// Items returns all items in slice
	// example:
	//
	//	slice := NewSlice[int]()
	//	slice.Appends(1, 2, 3, 4, 5)
	// then
	// slice.Items() -> [1, 2, 3, 4, 5]
	Items() []T

	// Get returns item by index
	// example:
	//
	//	slice := NewSlice[int]()
	//	slice.Appends(1, 2, 3, 4, 5)
	// then
	// slice.Get(0) -> 1
	Get(index int) T

	// Set sets item by index
	// example:
	//
	//	slice := NewSlice[int]()
	//	slice.Appends(1, 2, 3, 4, 5)
	//	slice.Set(0, 2)
	// then
	// slice.Items() -> [2, 2, 3, 4, 5]
	Set(index int, item T)

	// Length returns the length of slice
	// example:
	//
	//	slice := NewSlice[int]()
	//	slice.Appends(1, 2, 3, 4, 5)
	// then
	// slice.Length() -> 5
	Length() int

	// Capacity returns the capacity of slice
	// example:
	//
	//	slice := NewSlice[int]()
	//	slice.Appends(1, 2, 3, 4, 5)
	// then
	// slice.Capacity() -> 5
	Capacity() int
}

type slice[T any] struct {
	items []T
	mtx   sync.RWMutex
}

func (s *slice[T]) Append(item T) {
	s.mtx.Lock()
	s.items = append(s.items, item)
	s.mtx.Unlock()
}

func (s *slice[T]) Appends(items ...T) {
	if len(items) == 0 {
		return
	}

	s.mtx.Lock()
	s.items = append(s.items, items...)
	s.mtx.Unlock()
}

func (s *slice[T]) SubSlice(start, end int) (sub Slice[T]) {
	s.mtx.RLock()
	sub = NewSlice[T](s.items[start:end]...)
	s.mtx.RUnlock()

	return sub
}

func (s *slice[T]) Items() []T {
	s.mtx.RLock()
	items := make([]T, len(s.items))
	copy(items, s.items)
	s.mtx.RUnlock()

	return items
}

func (s *slice[T]) Get(index int) T {
	s.mtx.RLock()
	item := s.items[index]
	s.mtx.RUnlock()

	return item
}

func (s *slice[T]) Set(index int, item T) {
	s.mtx.Lock()
	s.items[index] = item
	s.mtx.Unlock()
}

func (s *slice[T]) Length() int {
	s.mtx.RLock()
	length := len(s.items)
	s.mtx.RUnlock()

	return length
}

func (s *slice[T]) Capacity() int {
	s.mtx.RLock()
	capacity := cap(s.items)
	s.mtx.RUnlock()

	return capacity
}

func NewSlice[T any](items ...T) Slice[T] {
	s := &slice[T]{
		items: make([]T, 0),
		mtx:   sync.RWMutex{},
	}

	if len(items) > 0 {
		s.items = append(s.items, items...)
	}

	return s
}
