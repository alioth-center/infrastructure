package exit

import (
	"embed"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/alioth-center/infrastructure/utils/concurrency"
)

var (
	// channel of blocking exit
	blockedChannel = make(chan struct{}, 1)

	// immediately exit if set to true
	exitImmediately = atomic.Bool{}

	//go:embed banner.txt
	banner embed.FS

	signalChannel = make(chan os.Signal, 1)
)

func init() {
	if os.Getenv("DISABLE_ALIOTH_BANNER_PRINT") == "" {
		bannerBytes, err := banner.ReadFile("banner.txt")
		if err == nil {
			// try to print banner, if error occurs, ignore it
			// you can delete this block if you don't want to print banner,
			// or you can change the banner.txt file to customize your banner
			fmt.Println(string(bannerBytes))
		}
	}

	exitImmediately.Store(true)
	eventList = concurrency.NewMap[string, EventHandler]()
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go handleNotify(signalChannel)
}

// handleNotify listens to the provided signal channel and manages the graceful shutdown
// of the application by executing registered exit functions. It waits for up to 10 seconds
// for these functions to complete before forcefully exiting.
//
// Parameters:
//
//	sg (chan os.Signal): The channel to receive OS signals indicating a termination request.
func handleNotify(sg chan os.Signal) {
	sig := <-sg

	fmt.Println("received signal:", sig.String(), "process will exit")
	fmt.Println("waiting for exit functions to finish...")

	wg := &sync.WaitGroup{}
	wg.Add(eventList.Length())

	go func() {
		// wait for 10 seconds, if exit functions are not finished, force exit
		time.Sleep(time.Second * 10)
		fmt.Println("exit functions are taking too long to finish, force exit")
		os.Exit(1)
	}()

	go eventList.Range(func(eventName string, function EventHandler) {
		// execute exit functions concurrently
		fmt.Println("start executing exit function:", eventName)
		go func(eventName string) {
			done := make(chan struct{})
			PrintlnUntilDone("executing exit function: "+eventName, time.Second, done)
			function(sig)
			fmt.Println("exit function executed:", eventName)
			close(done)
			wg.Done()
		}(eventName)
	})

	// wait for all exit functions to finish
	wg.Wait()

	// if call BlockedUntilTerminate, it will unblock
	if exitImmediately.Load() {
		os.Exit(0)
	}

	// otherwise, unblock the channel
	blockedChannel <- struct{}{}
}

// BlockedUntilTerminate blocks the current goroutine until a termination signal is received
// and all exit functions have completed. This function should be called to ensure the program
// does not exit immediately and waits for a proper shutdown sequence.
func BlockedUntilTerminate() {
	sync.OnceFunc(func() {
		exitImmediately.Store(false)
		<-blockedChannel
	})()
}

// Exit sends a termination signal to the signal channel, initiating the shutdown process.
func Exit() {
	signalChannel <- syscall.SIGTERM
}

// PrintlnUntilDone prints a message at regular intervals until the provided done channel is closed.
//
// Parameters:
//
//	message (string): The message to be printed periodically.
//	interval (time.Duration): The interval at which the message is printed.
//	done (chan struct{}): A channel that, when closed, stops the printing of the message.
func PrintlnUntilDone(message string, interval time.Duration, done chan struct{}) {
	go func() {
		for {
			select {
			case <-done:
				return
			case <-time.After(interval):
				fmt.Println(message)
			}
		}
	}()
}
