package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/alioth-center/infrastructure/cache"
	"github.com/alioth-center/infrastructure/exit"
	"github.com/alioth-center/infrastructure/utils/values"
	"github.com/go-redis/redis/v8"
)

const DriverName = "redis"

type Config struct {
	Address       string `json:"address,omitempty" yaml:"address,omitempty" xml:"address,omitempty"`
	Username      string `json:"username,omitempty" yaml:"username,omitempty" xml:"username,omitempty"`
	Password      string `json:"password,omitempty" yaml:"password,omitempty" xml:"password,omitempty"`
	DatabaseIndex int    `json:"database_index,omitempty" yaml:"database_index,omitempty" xml:"database_index,omitempty"`
	MaxRetries    int    `json:"max_retries,omitempty" yaml:"max_retries,omitempty" xml:"max_retries,omitempty"`
	TimeoutSecond int    `json:"timeout_second,omitempty" yaml:"timeout_second,omitempty" xml:"timeout_second,omitempty"`
	MaxLifeSecond int    `json:"max_life_second,omitempty" yaml:"max_life_second,omitempty" xml:"max_life_second,omitempty"`
	MaxOpen       int    `json:"max_open,omitempty" yaml:"max_open,omitempty" xml:"max_open,omitempty"`
	Prefix        string `json:"prefix,omitempty" yaml:"prefix,omitempty" xml:"prefix,omitempty"`
	KeySeparator  string `json:"key_separator,omitempty" yaml:"key_separator,omitempty" xml:"key_separator,omitempty"`
}

func newRedisClient(cfg Config) (rds *accessor, err error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Address,
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DatabaseIndex,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  time.Second * time.Duration(cfg.TimeoutSecond),
		ReadTimeout:  time.Second * time.Duration(cfg.TimeoutSecond),
		WriteTimeout: time.Second * time.Duration(cfg.TimeoutSecond),
		PoolSize:     cfg.MaxOpen,
		MaxConnAge:   time.Second * time.Duration(cfg.MaxLifeSecond),
	})

	_, pingErr := client.Ping(context.Background()).Result()
	if pingErr != nil {
		return values.Nil[*accessor](), fmt.Errorf("failed to connect redis server %s: %w", cfg.Address, pingErr)
	}

	// 初始化成功，需要注册退出函数
	exit.RegisterExitEvent(func(signal os.Signal) {
		_ = client.Close()
		fmt.Println("closed redis client")
	}, "CLOSE_REDIS_CONN")

	return &accessor{
		db: client,
		kb: keyBuilder{
			localRedisKeyPrefix: cfg.Prefix,
			redisKeySeparator:   cfg.KeySeparator,
		},
	}, nil
}

func NewRedisCache(cfg Config) (rds cache.Cache, err error) {
	return newRedisClient(cfg)
}

func NewRedisCounter(cfg Config) (rds cache.Counter, err error) {
	return newRedisClient(cfg)
}
