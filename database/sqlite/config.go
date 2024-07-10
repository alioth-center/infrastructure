package sqlite

import (
	"time"

	"github.com/alioth-center/infrastructure/database"
)

type Config struct {
	Database      string `yaml:"database,omitempty" json:"database,omitempty" xml:"database,omitempty"`
	Debug         bool   `yaml:"debug,omitempty" json:"debug,omitempty" xml:"debug,omitempty"`
	Stdout        string `yaml:"stdout,omitempty" json:"stdout,omitempty" xml:"stdout,omitempty"`
	Stderr        string `yaml:"stderr,omitempty" json:"stderr,omitempty" xml:"stderr,omitempty"`
	MaxIdle       int    `yaml:"max_idle,omitempty" json:"max_idle,omitempty" xml:"max_idle,omitempty"`
	MaxOpen       int    `yaml:"max_open,omitempty" json:"max_open,omitempty" xml:"max_open,omitempty"`
	MaxLifeSecond int    `yaml:"max_life_second,omitempty" json:"max_life_second,omitempty" xml:"max_life_second,omitempty"`
	TimeoutSecond int    `yaml:"timeout_second,omitempty" json:"timeout_second,omitempty" xml:"timeout_second,omitempty"`
}

func convertConfigToOptions(cfg Config) (opt database.Options) {
	return database.Options{
		DataSource: cfg.Database,
		DebugLog:   cfg.Debug,
		MaxIdle:    cfg.MaxIdle,
		MaxOpen:    cfg.MaxOpen,
		MaxLife:    time.Duration(cfg.MaxLifeSecond) * time.Second,
		Timeout:    time.Duration(cfg.TimeoutSecond) * time.Second,
	}
}
