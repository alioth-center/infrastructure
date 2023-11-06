package memory

import (
	"strconv"
	"testing"
	"time"
)

func TestMemoryCache(t *testing.T) {
	cache, _ := NewMemoryCache(Config{
		EnableInitiativeClean: true,
		CleanIntervalSecond:   1,
		MaxCleanMicroSecond:   10000,
		MaxCleanPercentage:    1,
	})

	for i := 0; i < 10086; i++ {
		go cache.StoreEX(nil, strconv.Itoa(i), "", time.Second+time.Duration(i)*time.Millisecond)
	}

	time.Sleep(time.Second * 5)
}

func BenchmarkMemoryCache(b *testing.B) {
	cache, _ := NewMemoryCache(Config{
		EnableInitiativeClean: true,
		CleanIntervalSecond:   1,
		MaxCleanMicroSecond:   100,
		MaxCleanPercentage:    10,
	})
	for i := 0; i < b.N; i++ {
		go cache.StoreEX(nil, strconv.Itoa(i), "", time.Second+time.Duration(i)*time.Millisecond)
	}
}
