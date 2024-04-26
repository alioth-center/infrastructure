package cls

import (
	"encoding/json"
	"fmt"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/utils/timezone"
	"strings"
)

const (
	TimeFormat = "2006.01.02-15:04:05.000-07:00"
)

type Logger struct {
	opts     Config
	cli      *client
	fallback logger.Logger
}

func (l *Logger) Debug(fields logger.Fields) {
	l.Log(logger.LevelDebug, fields)
}

func (l *Logger) Info(fields logger.Fields) {
	l.Log(logger.LevelInfo, fields)
}

func (l *Logger) Warn(fields logger.Fields) {
	l.Log(logger.LevelWarn, fields)
}

func (l *Logger) Error(fields logger.Fields) {
	l.Log(logger.LevelError, fields)
}

func (l *Logger) Fatal(fields logger.Fields) {
	l.Log(logger.LevelFatal, fields)
}

func (l *Logger) Panic(fields logger.Fields) {
	l.Log(logger.LevelPanic, fields)
}

func (l *Logger) Log(level logger.Level, fields logger.Fields) {
	if logger.LevelValueMap[l.opts.LogLevel] <= logger.LevelValueMap[level] {
		l.executeLog(fields.WithLevel(level))
	}

	if l.opts.LogLocal {
		l.fallback.Log(level, fields)
	}
}

func (l *Logger) Logf(level logger.Level, fields logger.Fields, format string, args ...any) {
	l.Log(level, fields.WithMessage(fmt.Sprintf(format, args...)))
}

func (l *Logger) Debugf(fields logger.Fields, format string, args ...any) {
	l.Debug(fields.WithMessage(fmt.Sprintf(format, args...)))
}

func (l *Logger) Infof(fields logger.Fields, format string, args ...any) {
	l.Info(fields.WithMessage(fmt.Sprintf(format, args...)))
}

func (l *Logger) Warnf(fields logger.Fields, format string, args ...any) {
	l.Warn(fields.WithMessage(fmt.Sprintf(format, args...)))
}

func (l *Logger) Errorf(fields logger.Fields, format string, args ...any) {
	l.Error(fields.WithMessage(fmt.Sprintf(format, args...)))
}

func (l *Logger) Fatalf(fields logger.Fields, format string, args ...any) {
	l.Fatal(fields.WithMessage(fmt.Sprintf(format, args...)))
}

func (l *Logger) Panicf(fields logger.Fields, format string, args ...any) {
	l.Panic(fields.WithMessage(fmt.Sprintf(format, args...)))
}

func (l *Logger) prepareField(field string) string {
	return strings.ReplaceAll(field, "-", "_")
}

func (l *Logger) prepareStructure(fields logger.Fields) map[string]string {
	entry := fields.Export()

	dataBytes := []byte("{}")
	if entry.Data != nil {
		bytes, marshalDataErr := json.Marshal(entry.Data)
		if marshalDataErr != nil {
			bytes = []byte("{}")
		}
		dataBytes = bytes
	}

	if entry.Extra == nil {
		entry.Extra = map[string]any{}
	}
	extraBytes, marshalExtraErr := json.Marshal(entry.Extra)
	if marshalExtraErr != nil {
		extraBytes = []byte("{}")
	}

	return map[string]string{
		"tid":      l.prepareField(entry.TraceID),
		"service":  l.prepareField(l.opts.Service),
		"instance": l.prepareField(l.opts.Instance),
		"time":     timezone.NowInTimezone(l.opts.Locale).Format(TimeFormat),
		"file":     entry.File,
		"level":    entry.Level,
		"func":     entry.Service,
		"message":  entry.Message,
		"data":     string(dataBytes),
		"extra":    string(extraBytes),
	}
}

func (l *Logger) executeLog(fields logger.Fields) {
	structure := l.prepareStructure(fields)
	l.cli.execute(structure)
}

func NewClsLogger(opts Config, fallback logger.Logger) (logger logger.Logger, err error) {
	cli, clsInitErr := newClsClient(opts, fallback)
	if clsInitErr != nil {
		return fallback, clsInitErr
	}

	logger = &Logger{
		opts:     opts,
		cli:      cli,
		fallback: fallback,
	}

	return logger, nil
}
