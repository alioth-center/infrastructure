package cache

import (
	"context"
	"time"
)

type CounterResultEnum int

const (
	CounterResultEnumFailed       CounterResultEnum = -1 // 操作失败
	CounterResultEnumNotEffective CounterResultEnum = 0  // 操作未改变任何内容
	CounterResultEnumSuccess      CounterResultEnum = 1  // 操作成功
)

func (e CounterResultEnum) GetValue() int {
	if e > 0 {
		return int(e)
	}

	return 0
}

type Counter interface {
	// Increase 将计数器的值增加delta，如果key不存在，则创建一个新的计数器，初始值为delta
	Increase(ctx context.Context, key string, delta uint64) (result CounterResultEnum)

	// IncreaseWithExpireWhenNotExist 将计数器的值增加delta；如果key不存在，则创建一个新的计数器，初始值为delta，过期时间为expire
	IncreaseWithExpireWhenNotExist(ctx context.Context, key string, delta uint64, expire time.Duration) (result CounterResultEnum)

	// SetExpire 设置计数器的过期时间，如果key已存在，则会将过期时间重置为expire
	SetExpire(ctx context.Context, key string, expire time.Duration) (result CounterResultEnum)

	// SetExpireWhenNotSet 将一个未设置过期时间的计数器的过期时间设置为expire；如果key不存在或已设置过期时间，则不会生效
	SetExpireWhenNotSet(ctx context.Context, key string, expire time.Duration) (result CounterResultEnum)

	// ExpireImmediately 计数器立即过期，如果key不存在，则不会生效
	ExpireImmediately(ctx context.Context, key string) (result CounterResultEnum)
}
