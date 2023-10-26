package exit

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
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

	wg := sync.WaitGroup{}
	for _, fn := range exitEvents {
		if fn != nil {
			wg.Add(1)
			go func(f EventHandler) {
				fn(signalString)
				wg.Done()
			}(fn)
		}
	}

	wg.Wait()
}
