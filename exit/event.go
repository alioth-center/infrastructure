package exit

import (
	"fmt"
)

type EventHandler func(sig string) string

type eventList struct {
	name    string
	handler EventHandler
}

var (
	exitEvents []eventList
)

// Register 向退出事件列表中注册事件，事件将在程序退出时执行，按注册顺序执行
func Register(fn EventHandler, name string) {
	if fn != nil {
		fmt.Println("register exit event", name)
		exitEvents = append(exitEvents, eventList{
			name:    name,
			handler: fn,
		})
	}
}
