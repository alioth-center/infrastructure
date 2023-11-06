package redis

import (
	"fmt"
	"strings"
)

var (
	globalRedisKeyPrefix string
)

func SetGlobalRedisKeyPrefix(key string) {
	globalRedisKeyPrefix = key
}

type keyBuilder struct {
	localRedisKeyPrefix string
	redisKeySeparator   string
}

func (kb keyBuilder) BuildKey(keys ...string) (result string) {
	builder := strings.Builder{}
	if globalRedisKeyPrefix != "" {
		builder.WriteString(globalRedisKeyPrefix)
	}

	if kb.localRedisKeyPrefix != "" {
		builder.WriteString(kb.redisKeySeparator)
		builder.WriteString(kb.localRedisKeyPrefix)
	}

	for _, key := range keys {
		if key != "" {
			builder.WriteString(kb.redisKeySeparator)
			builder.WriteString(key)
		}
	}

	return builder.String()
}

func (kb keyBuilder) BuildError(operation string, err error, keys ...string) (result error) {
	return fmt.Errorf("%s of key %s: %w", operation, kb.BuildKey(keys...), err)
}
