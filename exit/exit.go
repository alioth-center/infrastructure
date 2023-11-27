package exit

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	exc = make(chan struct{}, 1)
	set = false
)

func init() {
	sgChan := make(chan os.Signal, 1)
	signal.Notify(sgChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go handleNotify(sgChan)
}

func handleNotify(sg chan os.Signal) {
	sig, signalString := <-sg, ""
	switch sig {
	case syscall.SIGTERM:
		signalString = "SIGTERM"
	case syscall.SIGINT:
		signalString = "SIGINT"
	case syscall.SIGQUIT:
		signalString = "SIGQUIT"
	default:
		signalString = "SIG"
	}

	fmt.Printf("\n")
	fmt.Printf("waiting for %d exit functions to finish...\n", len(exitEvents))
	wg := &sync.WaitGroup{}
	for i, fn := range exitEvents {
		if fn.handler != nil {
			wg.Add(1)
			go func(i int, fn eventList) {
				defer wg.Done()
				fmt.Printf("executing %d/%d exit function %s with message: %s\n", i+1, len(exitEvents), fn.name, fn.handler(signalString))
			}(i, fn)
		}
	}

	wg.Wait()
	if !set {
		os.Exit(0)
	}
	exc <- struct{}{}
}

func BlockedUntilTerminate() {
	if set {
		return
	}

	set = true
	select {
	case <-exc:
		return
	}
}
