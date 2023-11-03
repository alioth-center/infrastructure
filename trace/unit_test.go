package trace

import (
	"fmt"
	"testing"
)

func TestStack(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			if i == 9 {
				t.Log("\n" + string(Stack(0)))
			}
		})
	}
}
