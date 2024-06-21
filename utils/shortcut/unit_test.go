package shortcut

import "testing"

func TestTernary(t *testing.T) {
	t.Run("TrueValue", func(t *testing.T) {
		if Ternary(true, 1, 0) != 1 {
			t.Error("Ternary(true, 1, 0) should return 1")
		}
	})

	t.Run("FalseValue", func(t *testing.T) {
		if Ternary(false, 1, 0) != 0 {
			t.Error("Ternary(false, 1, 0) should return 0")
		}
	})
}
