package cls

import (
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/utils/timezone"
)

type Config struct {
	Locale     timezone.TimeZone `json:"locale" yaml:"locale" xml:"locale"`
	Service    string            `json:"service" yaml:"service" xml:"service"`
	Instance   string            `json:"instance" yaml:"instance" xml:"instance"`
	SecretKey  string            `json:"secret_key" yaml:"secret_key" xml:"secret_key"`
	SecretID   string            `json:"secret_id" yaml:"secret_id" xml:"secret_id"`
	Endpoint   string            `json:"endpoint" yaml:"endpoint" xml:"endpoint"`
	TopicID    string            `json:"topic_id" yaml:"topic_id" xml:"topic_id"`
	MaxRetries int               `json:"max_retries" yaml:"max_retries" xml:"max_retries"`
	LogLocal   bool              `json:"log_local" yaml:"log_local" xml:"log_local"`
	LogLevel   logger.Level      `json:"log_level" yaml:"log_level" xml:"log_level"`
}
