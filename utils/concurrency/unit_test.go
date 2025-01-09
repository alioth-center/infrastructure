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
