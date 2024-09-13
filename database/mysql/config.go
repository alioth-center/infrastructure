package mysql

import (
	"fmt"
	"time"

	"github.com/alioth-center/infrastructure/database"
)

type Config struct {
	Server        string `yaml:"server,omitempty" json:"server,omitempty" xml:"server,omitempty"`
	Port          int    `yaml:"port,omitempty" json:"port,omitempty" xml:"port,omitempty"`
	Username      string `yaml:"username,omitempty" json:"username,omitempty" xml:"username,omitempty"`
	Password      string `yaml:"password,omitempty" json:"password,omitempty" xml:"password,omitempty"`
	Database      string `yaml:"database,omitempty" json:"database,omitempty" xml:"database,omitempty"`
	Charset       string `yaml:"charset,omitempty" json:"charset,omitempty" xml:"charset,omitempty"`
	Location      string `yaml:"location,omitempty" json:"location,omitempty" xml:"location,omitempty"`
	ParseTime     bool   `yaml:"parse_time,omitempty" json:"parse_time,omitempty" xml:"parse_time,omitempty"`
	Debug         bool   `yaml:"debug,omitempty" json:"debug,omitempty" xml:"debug,omitempty"`
	Stdout        string `yaml:"stdout,omitempty" json:"stdout,omitempty" xml:"stdout,omitempty"`
	Stderr        string `yaml:"stderr,omitempty" json:"stderr,omitempty" xml:"stderr,omitempty"`
	MaxIdle       int    `yaml:"max_idle,omitempty" json:"max_idle,omitempty" xml:"max_idle,omitempty"`
	MaxOpen       int    `yaml:"max_open,omitempty" json:"max_open,omitempty" xml:"max_open,omitempty"`
	MaxLifeSecond int    `yaml:"max_life_second,omitempty" json:"max_life_second,omitempty" xml:"max_life_second,omitempty"`
	TimeoutSecond int    `yaml:"timeout_second,omitempty" json:"timeout_second,omitempty" xml:"timeout_second,omitempty"`
}

func convertConfigToOptions(cfg Config) (opt database.Options) {
	parseTime := "True"
	if !cfg.ParseTime {
		parseTime = "False"
	}
	if cfg.Location == "" {
		cfg.Location = "Local"
	}
	if cfg.Charset == "" {
		cfg.Charset = "utf8mb4"
	}
	if cfg.Port == 0 {
		cfg.Port = 3306
	}

	// user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s", cfg.Username, cfg.Password, cfg.Server, cfg.Port, cfg.Database, cfg.Charset, parseTime, cfg.Location)
	return database.Options{
		DataSource: dsn,
		MaxIdle:    cfg.MaxIdle,
		MaxOpen:    cfg.MaxOpen,
		MaxLife:    time.Duration(cfg.MaxLifeSecond) * time.Second,
		Timeout:    time.Duration(cfg.TimeoutSecond) * time.Second,
	}
}
