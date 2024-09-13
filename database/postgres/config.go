package postgres

import (
	"fmt"
	"time"

	"github.com/alioth-center/infrastructure/database"
)

type Config struct {
	Host          string `yaml:"host,omitempty" json:"host,omitempty" xml:"host,omitempty"`
	Port          int    `yaml:"port,omitempty" json:"port,omitempty" xml:"port,omitempty"`
	Username      string `yaml:"username,omitempty" json:"username,omitempty" xml:"username,omitempty"`
	Password      string `yaml:"password,omitempty" json:"password,omitempty" xml:"password,omitempty"`
	Database      string `yaml:"database,omitempty" json:"database,omitempty" xml:"database,omitempty"`
	Charset       string `yaml:"charset,omitempty" json:"charset,omitempty" xml:"charset,omitempty"`
	Location      string `yaml:"location,omitempty" json:"location,omitempty" xml:"location,omitempty"`
	EnableSSL     bool   `yaml:"enable_ssl,omitempty" json:"enable_ssl,omitempty" xml:"enable_ssl,omitempty"`
	Debug         bool   `yaml:"debug,omitempty" json:"debug,omitempty" xml:"debug,omitempty"`
	Stdout        string `yaml:"stdout,omitempty" json:"stdout,omitempty" xml:"stdout,omitempty"`
	Stderr        string `yaml:"stderr,omitempty" json:"stderr,omitempty" xml:"stderr,omitempty"`
	MaxIdle       int    `yaml:"max_idle,omitempty" json:"max_idle,omitempty" xml:"max_idle,omitempty"`
	MaxOpen       int    `yaml:"max_open,omitempty" json:"max_open,omitempty" xml:"max_open,omitempty"`
	MaxLifeSecond int    `yaml:"max_life_second,omitempty" json:"max_life_second,omitempty" xml:"max_life_second,omitempty"`
	TimeoutSecond int    `yaml:"timeout_second,omitempty" json:"timeout_second,omitempty" xml:"timeout_second,omitempty"`
}

func convertConfigToOptions(cfg Config) (opt database.Options) {
	ssl := "disable"
	if cfg.EnableSSL {
		ssl = "enable"
	}
	if cfg.Location == "" {
		cfg.Location = "Asia/Shanghai"
	}

	// host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", cfg.Host, cfg.Username, cfg.Password, cfg.Database, cfg.Port, ssl, cfg.Location)
	return database.Options{
		DataSource: dsn,
		MaxIdle:    cfg.MaxIdle,
		MaxOpen:    cfg.MaxOpen,
		MaxLife:    time.Duration(cfg.MaxLifeSecond) * time.Second,
		Timeout:    time.Duration(cfg.TimeoutSecond) * time.Second,
	}
}
