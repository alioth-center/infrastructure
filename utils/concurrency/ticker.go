package concurrent

import (
	"sync"
	"sync/atomic"
)

var maintained = NewSlice[Resettable]()

type Resettable interface {
	Reset()
}

// ResetAll resets all instances in the maintained slice.
func ResetAll() {
	for _, r := range maintained.Items() {
		if r != nil {
			r.Reset()
		}
	}
}

// TickerInstance is an instance management interface based on a counter. It can periodically reset instances to avoid memory leaks, state accumulation, etc., and can also implement lazy loading within a certain range.
type TickerInstance[T any] interface {
	// Instance returns the current available instance. It must be called between TickStart() and TickEnd(), otherwise concurrency issues or null pointer issues may occur.
	Instance() T

	// TickStart starts an operation counting cycle. It must be called before TickEnd().
	TickStart()

	// TickEnd ends an operation counting cycle. It must be called after TickStart().
	TickEnd()

	// Reset immediately resets the instance to avoid long waits for the counting cycle to end.
	Reset()
}

// MaintainTickerInstance creates an instance manager based on a counter. It can periodically reset instances to avoid memory leaks, state accumulation, etc.
func MaintainTickerInstance[T any](ctor func() T, maxTick int) TickerInstance[T] {
	instance := &baseTickerInstance[T]{
		instance: ctor(),
		reset:    ctor,
		maxTick:  int32(maxTick),
	}

	maintained.Append(instance)
	return instance
}

type baseTickerInstance[T any] struct {
	instance T
	mtx      sync.RWMutex
	cnt      atomic.Int32
	reset    func() T
	maxTick  int32
	locked   atomic.Bool
}

func (b *baseTickerInstance[T]) TickStart() {
	b.mtx.RLock()
}

// TickEnd increments the tick counter and checks if it needs to reset the instance.
// If the counter exceeds the maximum tick value and the reset lock is successfully acquired,
// it performs the reset logic. Otherwise, it releases the read lock.
func (b *baseTickerInstance[T]) TickEnd() {
	b.cnt.Add(1)

	// Determine whether to reset, and execute the reset logic when the reset lock is successfully acquired
	if b.cnt.Load() > b.maxTick && b.locked.CompareAndSwap(false, true) {
		b.mtx.RUnlock()        // Unlock the read lock first to avoid deadlock when adding a write lock below
		b.mtx.Lock()           // Add a write lock to prevent other TickStart() from entering
		b.instance = b.reset() // Reset instance
		b.cnt.Store(0)         // Reset count
		b.locked.Store(false)  // Unlock reset status
		b.mtx.Unlock()         // Unlock the write lock to allow other TickStart() to enter
	} else {
		b.mtx.RUnlock()
	}
}

func (b *baseTickerInstance[T]) Instance() T {
	return b.instance
}

func (b *baseTickerInstance[T]) Reset() {
	b.mtx.Lock()
	b.instance = b.reset()
	b.cnt.Store(0)
	b.mtx.Unlock()
}
