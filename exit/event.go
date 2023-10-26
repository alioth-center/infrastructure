package exit

type EventHandler func(sig string)

var (
	exitEvents []EventHandler
)

// Register 向退出事件列表中注册事件，事件将在程序退出时执行，按注册顺序执行
func Register(fns ...EventHandler) {
	for _, fn := range fns {
		if fn == nil {
			continue
		}
		exitEvents = append(exitEvents, fn)
	}
}
