package exit

import (
	"fmt"
	"os"

	"github.com/alioth-center/infrastructure/trace"
	"github.com/alioth-center/infrastructure/utils/concurrency"
)

// eventList is a concurrent map that stores event names and their corresponding
// event handlers. It uses the concurrency.Map type to ensure thread-safe access.
var eventList concurrency.Map[string, EventHandler]

// EventHandler is a function type that takes an os.Signal as its parameter.
// This function is intended to handle the signal received during the program's
// exit process.
type EventHandler func(signal os.Signal)

// RegisterExitEvent registers an event handler for a specific exit event. The
// event handler is stored in the eventList map with the provided event name as
// the key. If the event handler function is not nil, it also prints details
// about the registration, including the event name, the location of the event
// handler function, and the location where the registration occurred.
//
// Parameters:
//
//	fn (EventHandler): The event handler function to be registered. It should
//	                   handle an os.Signal parameter.
//
//	eventName (string): The name of the event for which the handler is being
//	                    registered. This name is used as the key in the eventList map.
func RegisterExitEvent(fn EventHandler, eventName string) {
	if fn != nil {
		eventList.Set(eventName, fn)
		fmt.Println("exit event registered")
		fmt.Println("\t event name:", eventName)
		fmt.Println("\t event handler:", trace.FunctionLocation(fn))
		fmt.Println("\t registered at:", trace.Caller(0))
	}
}
