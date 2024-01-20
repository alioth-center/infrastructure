package redis

import (
	"os"
	"testing"
)

func TestRedisCache(t *testing.T) {
	if os.Getenv("ENABLE_REDIS_TEST") != "true" {
		t.Skip("skip redis test")
	}

	impl, initErr := NewRedisCache(Config{
		Address: "localhost:6379",
	})
	if initErr != nil {
		t.Fatal(initErr)
	}

	RunTestCase(t, impl)
}
