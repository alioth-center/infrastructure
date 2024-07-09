package exit

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func Fn(_ os.Signal) {
	time.Sleep(time.Second * 3)
	fmt.Println("exit")
}

func TestRegisterExitEvent(t *testing.T) {
	RegisterExitEvent(Fn, "testing")
}

func TestExitFunc(t *testing.T) {
	RegisterExitEvent(Fn, "testing")
	go func() {
		time.Sleep(time.Second)
		Exit()
	}()
	BlockedUntilTerminate()
}
