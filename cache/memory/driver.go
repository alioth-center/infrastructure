package memory

import (
	"github.com/alioth-center/infrastructure/cache"
	"github.com/alioth-center/infrastructure/exit"
	"sync"
	"time"
)

type Config struct {
	EnableInitiativeClean bool `json:"enable_initiative_clean,omitempty" yaml:"enable_initiative_clean,omitempty" xml:"enable_initiative_clean,omitempty"`
	CleanIntervalSecond   int  `json:"clean_interval_second,omitempty" yaml:"clean_interval_second,omitempty" xml:"clean_interval_second,omitempty"`
	MaxCleanMicroSecond   int  `json:"max_clean_micro_second,omitempty" yaml:"max_clean_micro_second,omitempty" xml:"max_clean_micro_second,omitempty"`
	MaxCleanPercentage    int  `json:"max_clean_percentage,omitempty" yaml:"max_clean_percentage,omitempty" xml:"max_clean_percentage,omitempty"`
}

func NewMemoryCache(cfg Config) (mc cache.Cache, err error) {
	memoryCache := &accessor{
		mtx: sync.RWMutex{},
		db:  map[string]entry{},
		ec:  make(chan struct{}, 1),
	}

	if cfg.EnableInitiativeClean {
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
