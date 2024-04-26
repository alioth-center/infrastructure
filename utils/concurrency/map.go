package concurrency

import (
	"bytes"
	"encoding/binary"
	"hash/fnv"
	"modernc.org/mathutil"
	"sync"
)

// Map is a thread-safe map
type Map[K comparable, V any] interface {
	// Get returns value by key
	// example:
	//
	//	m := NewMap[int, string]()
	//	m.Set(1, "test")
	// then
	// m.Get(1) -> "test", true
	Get(key K) (V, bool)

	// Set sets value by key
	// example:
	//
	//	m := NewMap[int, string]()
	//	m.Set(1, "test")
	// then
	// m.Get(1) -> "test", true
	Set(key K, value V)

	// Delete deletes value by key
	// example:
	//
	//	m := NewMap[int, string]()
	//	m.Set(1, "test")
	//	m.Delete(1)
	// then
	// m.Get(1) -> "", false
	Delete(key K)

	// Keys returns all keys in map
	// example:
	//
	//	m := NewMap[int, string]()
	//	m.Set(1, "test")
	//	m.Set(2, "test2")
	// then
	// m.Keys() -> [1, 2]
	Keys() []K

	// Values returns all values in map
	// example:
	//
	//	m := NewMap[int, string]()
	//	m.Set(1, "test")
	//	m.Set(2, "test2")
	// then
	// m.Values() -> ["test", "test2"]
	Values() []V

	// Range iterates all items in map
	// example:
	//
	//	m := NewMap[int, string]()
	//	m.Set(1, "test")
	//	m.Set(2, "test2")
	//	m.Range(func(key int, value string) {
	//		fmt.Println(key, value)
	//	})
	// then
	// 1 test
	// 2 test2
	Range(f func(key K, value V))

	// Origin returns a copy of origin map
	// example:
	//
	//	m := NewMap[int, string]()
	//	m.Set(1, "test")
	//	m.Set(2, "test2")
	//	o := m.Origin()
	// then
	//  o -> {1: "test", 2: "test2"}
	Origin() map[K]V

	// Length returns the length of map
	// example:
	//
	//	m := NewMap[int, string]()
	//	m.Set(1, "test")
	//	m.Set(2, "test2")
	// then
	// m.Length() -> 2
	Length() int

	// Clear clears all items in map
	// example:
	//
	//	m := NewMap[int, string]()
	//	m.Set(1, "test")
	//	m.Set(2, "test2")
	//	m.Clear()
	// then
	// m.Length() -> 0
	Clear()
}

type threadSafeMap[K comparable, V any] struct {
	items map[K]V
	mtx   sync.RWMutex
}

func (t *threadSafeMap[K, V]) Get(key K) (V, bool) {
	t.mtx.RLock()
	value, exist := t.items[key]
	t.mtx.RUnlock()

	return value, exist
}

func (t *threadSafeMap[K, V]) Set(key K, value V) {
	t.mtx.Lock()
	t.items[key] = value
	t.mtx.Unlock()
}

func (t *threadSafeMap[K, V]) Delete(key K) {
	t.mtx.Lock()
	delete(t.items, key)
	t.mtx.Unlock()
}

func (t *threadSafeMap[K, V]) Keys() []K {
	t.mtx.RLock()
	keys := make([]K, 0, len(t.items))
	for key := range t.items {
		keys = append(keys, key)
	}
	t.mtx.RUnlock()

	return keys
}

func (t *threadSafeMap[K, V]) Values() []V {
	t.mtx.RLock()
	values := make([]V, 0, len(t.items))
	for _, value := range t.items {
		values = append(values, value)
	}
	t.mtx.RUnlock()

	return values
}

func (t *threadSafeMap[K, V]) Range(f func(key K, value V)) {
	t.mtx.RLock()
	for key, value := range t.items {
		f(key, value)
	}
	t.mtx.RUnlock()
}

func (t *threadSafeMap[K, V]) Origin() map[K]V {
	t.mtx.RLock()
	items := make(map[K]V, len(t.items))
	for key, value := range t.items {
		items[key] = value
	}
	t.mtx.RUnlock()

	return items
}

func (t *threadSafeMap[K, V]) Length() int {
	t.mtx.RLock()
	length := len(t.items)
	t.mtx.RUnlock()

	return length
}

func (t *threadSafeMap[K, V]) Clear() {
	t.mtx.Lock()
	t.items = make(map[K]V)
	t.mtx.Unlock()
}

// NewMap returns a new thread-safe map
func NewMap[K comparable, V any]() Map[K, V] {
	return &threadSafeMap[K, V]{
		mtx:   sync.RWMutex{},
		items: make(map[K]V),
	}
}

// hashFunction is a hash function for hashMap, use hash/fnv to hash key
func hashFunction[T comparable](key T) (hash uint64) {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, key) // nolint:errcheck

	h := fnv.New64a()
	_, _ = h.Write(buf.Bytes()) // nolint:errcheck

	return h.Sum64()
}

// hashMap is a thread-safe map with hash function
type hashMap[K comparable, V any] struct {
	maxNodes uint64
	mtx      sync.RWMutex
	items    []Map[K, V]
}

func (h *hashMap[K, V]) index(key K) int {
	return int(hashFunction(key) % h.maxNodes)
}

func (h *hashMap[K, V]) Get(key K) (V, bool) {
	h.mtx.RLock()
	value, exist := h.items[h.index(key)].Get(key)
	h.mtx.RUnlock()

	return value, exist
}

func (h *hashMap[K, V]) Set(key K, value V) {
	h.mtx.RLock()
	h.items[h.index(key)].Set(key, value)
	h.mtx.RUnlock()
}

func (h *hashMap[K, V]) Delete(key K) {
	h.mtx.RLock()
	h.items[h.index(key)].Delete(key)
	h.mtx.RUnlock()
}

func (h *hashMap[K, V]) Keys() []K {
	h.mtx.RLock()
	results := make([]K, 0)
	for _, item := range h.items {
		results = append(results, item.Keys()...)
	}
	h.mtx.RUnlock()

	return results
}

func (h *hashMap[K, V]) Values() []V {
	h.mtx.RLock()
	results := make([]V, 0)
	for _, item := range h.items {
		results = append(results, item.Values()...)
	}
	h.mtx.RUnlock()

	return results
}

func (h *hashMap[K, V]) Range(f func(key K, value V)) {
	h.mtx.RLock()
	for _, item := range h.items {
		item.Range(f)
	}
	h.mtx.RUnlock()
}

func (h *hashMap[K, V]) Origin() map[K]V {
	h.mtx.RLock()
	results := make(map[K]V)
	for _, item := range h.items {
		for key, value := range item.Origin() {
			results[key] = value
		}
	}
	h.mtx.RUnlock()

	return results
}

func (h *hashMap[K, V]) Length() int {
	h.mtx.RLock()
	length := 0
	for _, item := range h.items {
		length += item.Length()
	}
	h.mtx.RUnlock()

	return length
}

func (h *hashMap[K, V]) Clear() {
	h.mtx.Lock()
	for _, item := range h.items {
		item.Clear()
	}
	h.mtx.Unlock()
}

type HashMapNodeOption uint64

const (
	HashMapNodeOptionSmallSize  HashMapNodeOption = 7
	HashMapNodeOptionMediumSize HashMapNodeOption = 47
	HashMapNodeOptionLargeSize  HashMapNodeOption = 97
	HashMapNodeOptionExtraSize  HashMapNodeOption = 997
	HashMapNodeOptionHugeSize   HashMapNodeOption = 9973
)

// CustomHashMapNodeOptions returns a custom HashMapNodeOption
func CustomHashMapNodeOptions(want uint64) HashMapNodeOption {
	if HashMapNodeOption(want) < HashMapNodeOptionSmallSize {
		return HashMapNodeOptionSmallSize
	}
	if HashMapNodeOption(want) < HashMapNodeOptionMediumSize {
		return HashMapNodeOptionMediumSize
	}

	for i := HashMapNodeOption(want); i > HashMapNodeOptionMediumSize; i-- {
		if mathutil.IsPrimeUint64(uint64(i)) {
			return i
		}
	}

	return HashMapNodeOptionMediumSize
}

// NewHashMap returns a new thread-safe map with hash function
//
// Attention: the performance of hashMap is much better than Map, but this structure CANNOT be used for persistence, the hash value of the key will change after restarting the program
func NewHashMap[K comparable, V any](maxNodes HashMapNodeOption) Map[K, V] {
	var realNodes uint64
	switch maxNodes {
	case HashMapNodeOptionSmallSize:
		realNodes = uint64(HashMapNodeOptionSmallSize)
	case HashMapNodeOptionMediumSize:
		realNodes = uint64(HashMapNodeOptionMediumSize)
	case HashMapNodeOptionLargeSize:
		realNodes = uint64(HashMapNodeOptionLargeSize)
	case HashMapNodeOptionExtraSize:
		realNodes = uint64(HashMapNodeOptionExtraSize)
	case HashMapNodeOptionHugeSize:
		realNodes = uint64(HashMapNodeOptionHugeSize)
	default:
		realNodes = uint64(CustomHashMapNodeOptions(uint64(maxNodes)))
	}

	items := make([]Map[K, V], maxNodes)
	for i := uint64(0); i < realNodes; i++ {
		items[i] = NewMap[K, V]()
	}

	return &hashMap[K, V]{
		maxNodes: realNodes,
		mtx:      sync.RWMutex{},
		items:    items,
	}
}
