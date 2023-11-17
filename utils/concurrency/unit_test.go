package concurrency

import (
	"math/rand"
	"testing"
	"time"
)

func TestConcurrency(t *testing.T) {
	st := time.Now()
	var promises = make([]Promise[string], 10)
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
