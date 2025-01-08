package concurrency

import (
	"errors"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestConcurrency(t *testing.T) {
	st := time.Now()
	promises := make([]Promise[string], 10)
	for i := 0; i < 10; i++ {
		fn := func() string {
			time.Sleep(1 * time.Second)
			if rand.Int()%2 == 1 {
				panic("random panic")
			}
			return "hello"
		}

		promises[i] = Async(fn)
	}

	for _, promise := range promises {
		t.Log(Await(promise))
	}

	t.Log("cost:", time.Since(st))

	p := Async(func() string {
		time.Sleep(1 * time.Second)
		return "hello"
	})

	t.Log(Await(p))
	t.Log(Await(p))
}

func TestRecoverErr(t *testing.T) {
	t.Run("RecoverErr:String", func(t *testing.T) {
		defer func() {
			if err := RecoverErr(recover()); err != nil {
				t.Log(err)
			}
		}()

		panic("test")
	})

	t.Run("RecoverErr:Error", func(t *testing.T) {
		defer func() {
			if err := RecoverErr(recover()); err != nil {
				t.Log(err)
			}
		}()

		panic(errors.New("test"))
	})

	t.Run("RecoverErr:Nil", func(t *testing.T) {
		defer func() {
			if err := RecoverErr(recover()); err != nil {
				t.Log(err)
			}
		}()

		panic(nil)
	})
}

func TestSlice(t *testing.T) {
	t.Run("Slice:Append", func(t *testing.T) {
		s := NewSlice[int]()
		s.Append(1)
		s.Append(2)

		if s.Length() != 2 {
			t.Error("slice length not match")
		}
		if s.Capacity() != 2 {
			t.Error("slice capacity not match")
		}
		for i, v := range s.Items() {
			if v != i+1 {
				t.Error("slice item not match")
			}
		}
	})

	t.Run("Slice:Appends", func(t *testing.T) {
		s := NewSlice[int]()
		s.Appends(1, 2, 3, 4, 5)

		if s.Length() != 5 {
			t.Error("slice length not match")
		}
		if s.Capacity() < 5 {
			t.Error("slice capacity not match")
		}
		for i, v := range s.Items() {
			if v != i+1 {
				t.Error("slice item not match")
			}
		}
	})

	t.Run("Slice:SubSlice", func(t *testing.T) {
		s := NewSlice[int]()
		s.Appends(1, 2, 3, 4, 5)

		sub := s.SubSlice(1, 3)
		if sub.Length() != 2 {
			t.Error("slice length not match")
		}
		if sub.Capacity() != 2 {
			t.Error("slice capacity not match")
		}
		for i, v := range sub.Items() {
			if v != i+2 {
				t.Error("slice item not match")
			}
		}
	})

	t.Run("Slice:Items", func(t *testing.T) {
		s := NewSlice[int]()
		s.Appends(1, 2, 3, 4, 5)

		for i, v := range s.Items() {
			if v != i+1 {
				t.Error("slice item not match")
			}
		}
	})

	t.Run("Slice:Get", func(t *testing.T) {
		s := NewSlice[int]()
		s.Appends(1, 2, 3, 4, 5)

		for i := 0; i < s.Length(); i++ {
			if s.Get(i) != i+1 {
				t.Error("slice item not match")
			}
		}
	})

	t.Run("Slice:Set", func(t *testing.T) {
		s := NewSlice[int]()
		s.Appends(1, 2, 3, 4, 5)

		s.Set(0, 2)
		if s.Get(0) != 2 {
			t.Error("slice item not match")
		}
	})

	t.Run("Slice:Length", func(t *testing.T) {
		s := NewSlice[int]()
		s.Appends(1, 2, 3, 4, 5)

		if s.Length() != 5 {
			t.Error("slice length not match")
		}
	})

	t.Run("Slice:Capacity", func(t *testing.T) {
		s := NewSlice[int]()
		s.Appends(1, 2, 3, 4, 5)

		if s.Capacity() != 6 {
			t.Error("slice capacity not match")
		}
	})

	t.Run("Slice:Concurrent", func(t *testing.T) {
		s := NewSlice[int]()
		wg := sync.WaitGroup{}
		wg.Add(10000)
		for i := 0; i < 100; i++ {
			go func() {
				for i := 0; i < 100; i++ {
					s.Append(1)
					wg.Done()
				}
			}()
		}

		wg.Wait()
		if s.Length() != 10000 {
			t.Error("slice length not match")
		}
		if s.Capacity() < 10000 {
			t.Error("slice capacity not match")
		}
	})
}

func TestMap(t *testing.T) {
	testCases := []struct {
		CaseName string
		CaseFn   func(m Map[int, string], t *testing.T)
	}{
		{
			CaseName: "Map:Set",
			CaseFn: func(m Map[int, string], t *testing.T) {
				m.Set(1, "test")
				m.Set(2, "test2")

				if m.Length() != 2 {
					t.Error("map length not match")
				}
				if value, exist := m.Get(1); value != "test" || !exist {
					t.Error("map item not match")
				}
				if value, exist := m.Get(2); value != "test2" || !exist {
					t.Error("map item not match")
				}
			},
		},
		{
			CaseName: "Map:Get",
			CaseFn: func(m Map[int, string], t *testing.T) {
				m.Set(1, "test")
				m.Set(2, "test2")

				if m.Length() != 2 {
					t.Error("map length not match")
				}
				if value, exist := m.Get(1); value != "test" || !exist {
					t.Error("map item not match")
				}
				if value, exist := m.Get(2); value != "test2" || !exist {
					t.Error("map item not match")
				}
			},
		},
		{
			CaseName: "Map:Delete",
			CaseFn: func(m Map[int, string], t *testing.T) {
				m.Set(1, "test")
				m.Set(2, "test2")
				m.Delete(1)

				if m.Length() != 1 {
					t.Error("map length not match")
				}
				if value, exist := m.Get(1); value != "" || exist {
					t.Error("map item not match")
				}
				if value, exist := m.Get(2); value != "test2" || !exist {
					t.Error("map item not match")
				}
			},
		},
		{
			CaseName: "Map:Keys",
			CaseFn: func(m Map[int, string], t *testing.T) {
				m.Set(1, "test")
				m.Set(2, "test2")

				keys := m.Keys()
				if len(keys) != 2 {
					t.Error("map keys length not match")
				}
				if keys[0] != 1 && keys[1] != 1 {
					t.Error("map keys not match")
				}
				if keys[0] != 2 && keys[1] != 2 {
					t.Error("map keys not match")
				}
			},
		},
		{
			CaseName: "Map:Values",
			CaseFn: func(m Map[int, string], t *testing.T) {
				m.Set(1, "test")
				m.Set(2, "test2")

				values := m.Values()
				if len(values) != 2 {
					t.Error("map values length not match")
				}
				if values[0] != "test" && values[1] != "test" {
					t.Error("map values not match")
				}
				if values[0] != "test2" && values[1] != "test2" {
					t.Error("map values not match")
				}
			},
		},
		{
			CaseName: "Map:Range",
			CaseFn: func(m Map[int, string], t *testing.T) {
				m.Set(1, "test")
				m.Set(2, "test2")

				m.Range(func(key int, value string) {
					if key == 1 && value != "test" {
						t.Error("map item not match")
					}
					if key == 2 && value != "test2" {
						t.Error("map item not match")
					}
				})
			},
		},
		{
			CaseName: "Map:Origin",
			CaseFn: func(m Map[int, string], t *testing.T) {
				m.Set(1, "test")
				m.Set(2, "test2")

				origin := m.Origin()
				if len(origin) != 2 {
					t.Error("map origin length not match")
				}
				if origin[1] != "test" || origin[2] != "test2" {
					t.Error("map origin not match")
				}
			},
		},
		{
			CaseName: "Map:Length",
			CaseFn: func(m Map[int, string], t *testing.T) {
				m.Set(1, "test")
				m.Set(2, "test2")

				if m.Length() != 2 {
					t.Error("map length not match")
				}
			},
		},
		{
			CaseName: "Map:Clear",
			CaseFn: func(m Map[int, string], t *testing.T) {
				m.Set(1, "test")
				m.Set(2, "test2")
				m.Clear()

				if m.Length() != 0 {
					t.Error("map length not match")
				}
			},
		},
		{
			CaseName: "Map:Concurrent",
			CaseFn: func(m Map[int, string], t *testing.T) {
				wg := sync.WaitGroup{}
				wg.Add(10000)
				for i := 0; i < 100; i++ {
					go func() {
						for i := 0; i < 100; i++ {
							m.Set(i, "test")
							wg.Done()
						}
					}()
				}

				wg.Wait()
				if m.Length() != 100 {
					t.Error("map length not match", m.Length())
				}
			},
		},
	}

	t.Run("Map:ThreadSafe", func(t *testing.T) {
		for _, testCase := range testCases {
			t.Run(testCase.CaseName, func(t *testing.T) {
				m := NewMap[int, string]()
				testCase.CaseFn(m, t)
			})
		}
	})

	t.Run("Map:HashMap", func(t *testing.T) {
		for _, testCase := range testCases {
			t.Run(testCase.CaseName, func(t *testing.T) {
				m := NewHashMap[int, string](HashMapNodeOptionSmallSize)
				testCase.CaseFn(m, t)
			})
		}
	})
}

func TestChain_Continue(t *testing.T) {
	t.Run("all handlers succeed", func(t *testing.T) {
		d := 123
		handler := func(ctx context.Context, data int) error {
			d++
			return nil
		}
		chain := NewChain(handler, handler, handler, handler)
		chain.Continue(context.Background(), 123)

		if err := chain.Error(); err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if d != 127 {
			t.Errorf("expected data to be 127, got %d", d)
		}
	})

	t.Run("context canceled", func(t *testing.T) {
		d := 123
		handler := func(ctx context.Context, data int) error {
			d++
			return nil
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Immediately cancel the context

		chain := NewChain(handler, handler)
		chain.Continue(ctx, 123)

		if err := chain.Error(); !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled error, got %v", err)
		}
		if d != 123 {
			t.Errorf("expected data to be 123, got %d", d)
		}
	})

	t.Run("handler returns error", func(t *testing.T) {
		expectedErr := errors.New("handler error")
		handler1 := func(ctx context.Context, data int) error {
			return nil
		}
		handler2 := func(ctx context.Context, data int) error {
			return expectedErr
		}
		handler3 := func(ctx context.Context, data int) error {
			t.Errorf("handler3 should not be executed")
			return nil
		}

		chain := NewChain(handler1, handler2, handler3)
		chain.Continue(context.Background(), 123)

		if err := chain.Error(); !errors.Is(err, expectedErr) {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
	})
}

func TestChain_ConcurrentExecution(t *testing.T) {
	t.Run("concurrent execution", func(t *testing.T) {
		handler := func(ctx context.Context, data int) error {
			time.Sleep(100 * time.Millisecond) // Simulate work
			return nil
		}

		chain := NewChain(handler, handler, handler)
		ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
		defer cancel()

		go chain.Continue(ctx, 123)
		chain.Continue(ctx, 123)

		if err := chain.Error(); err != nil && !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("expected deadline exceeded or nil, got %v", err)
		}
	})
}

type TestInstance struct {
	value int
}

func TestTickerInstance_BasicFunctionality(t *testing.T) {
	resetFunc := func() TestInstance {
		return TestInstance{value: 0}
	}
	maxTickCount := int32(5)
	instance := &baseTickerInstance[TestInstance]{
		instance: TestInstance{value: 1},
		reset:    resetFunc,
		maxTick:  maxTickCount,
	}

	for i := 0; i < 10; i++ {
		instance.TickStart()
		v := instance.Instance().value
		instance.TickEnd()

		if v != 1 && i < int(maxTickCount) {
			// For the first 5 ticks, value should be 1 before reset
			t.Fatalf("Expected instance value to be 1, got %d", v)
		} else if v != 0 && i > int(maxTickCount) {
			// After reset (after maxTickCount ticks), value should be reset to 0
			t.Fatalf("Expected instance value to be 0 after reset, got %d", v)
		}
	}
}

func TestTickerInstance_ConcurrentAccess(t *testing.T) {
	var resetCount atomic.Int32
	resetFunc := func() TestInstance {
		resetCount.Add(1)
		return TestInstance{value: 0}
	}
	maxTickCount := int32(4)
	instance := &baseTickerInstance[TestInstance]{
		instance: TestInstance{value: 1},
		reset:    resetFunc,
		maxTick:  maxTickCount,
	}

	t.Run("NonConcurrentAccess", func(t *testing.T) {
		// 非并发访问下，重置次数应该严格等于期望的次数
		var wg sync.WaitGroup
		const numGoroutines = 1
		const numOpsPerGoroutine = 10000

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < numOpsPerGoroutine; j++ {
					instance.TickStart()
					instance.Instance()
					instance.TickEnd()
				}
			}()
		}

		wg.Wait()
		expectedResets := int32((numGoroutines * numOpsPerGoroutine) / (int(maxTickCount) + 1))
		actualResets := resetCount.Load()

		if actualResets != expectedResets {
			t.Fatalf("Expected at least %d resets, got %d", expectedResets, actualResets)
		}
	})

	t.Run("ConcurrentAccess", func(t *testing.T) {
		// 并发访问下，因为重置的写锁要等正常操作的读锁释放，在实际重置实例的时候，一般实际 tick 次数会高于限制的值，所以重置次数会低于期望的次数
		resetCount.Store(0)
		var wg sync.WaitGroup
		const numGoroutines = 10
		const numOpsPerGoroutine = 10000

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < numOpsPerGoroutine; j++ {
					instance.TickStart()
					instance.Instance()
					instance.TickEnd()
				}
			}()
		}

		wg.Wait()
		expectedResets := int32((numGoroutines * numOpsPerGoroutine) / (int(maxTickCount) + 1))
		actualResets := resetCount.Load()

		if actualResets > expectedResets {
			t.Fatalf("Expected at most %d resets, got %d", expectedResets, actualResets)
		}
	})

}