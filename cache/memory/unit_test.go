package memory

import (
	"strconv"
	"testing"
	"time"
)

func TestMemoryCache(t *testing.T) {
	impl := NewMemoryCache(Config{
		EnableInitiativeClean: true,
		CleanIntervalSecond:   1,
		MaxCleanMicroSecond:   100,
		MaxCleanPercentage:    10,
	})

	RunCacheTestCases(t, impl)
}

func TestMemoryCounter(t *testing.T) {
	impl := NewMemoryCounter(Config{
		EnableInitiativeClean: true,
		CleanIntervalSecond:   1,
		MaxCleanMicroSecond:   100,
		MaxCleanPercentage:    10,
	})

	RunCounterTestCases(t, impl)
}

func BenchmarkMemoryCache(b *testing.B) {
	cache := NewMemoryCache(Config{
		EnableInitiativeClean: true,
		CleanIntervalSecond:   1,
		MaxCleanMicroSecond:   100,
		MaxCleanPercentage:    10,
	})
	for i := 0; i < b.N; i++ {
		go cache.StoreEX(nil, strconv.Itoa(i), "", time.Second+time.Duration(i)*time.Millisecond)
	}
}
