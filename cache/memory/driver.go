package memory

import (
	"github.com/alioth-center/infrastructure/cache"
	"github.com/alioth-center/infrastructure/exit"
	"sync"
	"time"
)

type Config struct {
	EnableClean         bool
	CleanIntervalSecond int
	MaxCleanMicroSecond int
	MaxCleanPercentage  int
}

func NewMemoryCache(cfg Config) (mc cache.Cache, err error) {
	memoryCache := &accessor{
		mtx: sync.RWMutex{},
		db:  map[string]entry{},
		ec:  make(chan struct{}, 1),
	}

	if cfg.EnableClean {
		interval, maxExec := time.Second*time.Duration(cfg.CleanIntervalSecond), time.Microsecond*time.Duration(cfg.MaxCleanMicroSecond)
		go memoryCache.cleanCache(interval, maxExec, cfg.MaxCleanPercentage)

		// 启动了主动淘汰策略，需要注册退出事件
		exit.Register(func(_ string) string {
			memoryCache.close()
			return "closed memory cache"
		}, "memory cache")
	}

	return memoryCache, nil
}
