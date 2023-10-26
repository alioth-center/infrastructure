package exit

import (
	"sync"
	"testing"
	"time"
)

func TestExit(t *testing.T) {
	idx, wg := 0, sync.WaitGroup{}
	wg.Add(3)
	exitFunc := func() func(sig string) {
		return func(sig string) {
			// do something
			time.Sleep(time.Second * time.Duration(idx*5))
			println(idx, sig, "exit")
			wg.Done()
			idx++
		}
	}
	Register(exitFunc(), exitFunc(), exitFunc())

	// press ctrl + c in terminal or click stop button in goland

	wg.Wait()
	t.Logf("exit success if you see this log")
}
