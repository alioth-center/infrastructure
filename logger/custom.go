package logger

import (
	"fmt"
	"os"
	"strings"
)

const (
	serviceEnvKey  = "AC_SERVICE"
	extraFieldsKey = "AC_EXTRA_FIELDS"
)

// NewCustomLoggerWithOpts creates and returns a new custom logger instance
// with the specified options. This function initializes a custom logger
// with default values and then applies the provided options to it.
//
// The logger is configured with default options for standard writer, JSON
// formatting, and log level set to Info. Additional options can be passed
// and will be applied in the order they are provided.
//
// Parameters:
//
//	opts (...Option): A variadic parameter that allows passing multiple
//	                  Option functions to customize the logger.
//
// Returns:
//
//	Logger: A configured logger instance with applied options.
func NewCustomLoggerWithOpts(opts ...Option) Logger {
	c := &customLogger{
		hooks:      map[Level][]func(Fields){},
		level:      LevelDebug,
		marshaller: defaultMarshaller,
	}

	// Apply default options before any user-provided options
	opts = append([]Option{
		WithStdWriterOpts(),
		WithJsonFormatOpts(),
		WithLevelOpts(LevelInfo),
	}, opts...)

	// Inject service field
	attachField := NewFields()
	if srv := os.Getenv(serviceEnvKey); srv != "" {
		attachField = attachField.WithService(srv)
	}

	// Inject extra fields
	if extraKeys := strings.Split(strings.TrimSpace(os.Getenv(extraFieldsKey)), ","); len(extraKeys) > 0 {
		for _, key := range extraKeys {
			if value := os.Getenv(key); value != "" {
				attachField = attachField.WithField(key, value)
			}
		}
	}

	// If inject fields not nil, inject it
	if entry := attachField.Export(); len(entry.Extra) != 0 || entry.Service != "" {
		opts = append(opts, WithAttachFields(attachField))
	}

	// Apply all options to the custom logger
	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}

	// Add a default hook to write logs using the configured marshaller
	hook := func(fields Fields) { c.writer.Write(c.marshaller(fields)) }
	WithHookOpts(hook, LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal, LevelPanic)(c)

	return c
}

type customLogger struct {
	hooks      map[Level][]func(Fields)
	level      Level
	marshaller func(Fields) []byte
	writer     Writer
	attach     Fields
}

func (c customLogger) Debug(fields Fields) {
	c.log(LevelDebug, fields)
}

func (c customLogger) Info(fields Fields) {
	c.log(LevelInfo, fields)
}

func (c customLogger) Warn(fields Fields) {
	c.log(LevelWarn, fields)
}

func (c customLogger) Error(fields Fields) {
	c.log(LevelError, fields)
}

func (c customLogger) Fatal(fields Fields) {
	c.log(LevelFatal, fields)
}

func (c customLogger) Panic(fields Fields) {
	c.log(LevelPanic, fields)
}

func (c customLogger) Log(level Level, fields Fields) {
	c.log(level, fields)
}

func (c customLogger) Logf(level Level, fields Fields, format string, args ...any) {
	c.log(level, fields.WithMessage(fmt.Sprintf(format, args...)))
}

func (c customLogger) Debugf(fields Fields, format string, args ...any) {
	c.log(LevelDebug, fields.WithMessage(fmt.Sprintf(format, args...)))
}

func (c customLogger) Infof(fields Fields, format string, args ...any) {
	c.log(LevelInfo, fields.WithMessage(fmt.Sprintf(format, args...)))
}

func (c customLogger) Warnf(fields Fields, format string, args ...any) {
	c.log(LevelWarn, fields.WithMessage(fmt.Sprintf(format, args...)))
}

func (c customLogger) Errorf(fields Fields, format string, args ...any) {
	c.Log(LevelError, fields.WithMessage(fmt.Sprintf(format, args...)))
}

func (c customLogger) Fatalf(fields Fields, format string, args ...any) {
	c.Log(LevelFatal, fields.WithMessage(fmt.Sprintf(format, args...)))
}

func (c customLogger) Panicf(fields Fields, format string, args ...any) {
	c.log(LevelPanic, fields.WithMessage(fmt.Sprintf(format, args...)))
}

func (c customLogger) log(level Level, fields Fields) {
	if c.level.shouldLog(level) {
		callbacks := c.hooks[level]
		if c.attach != nil {
			fields = fields.WithAttachFields(c.attach)
		}
		for _, callback := range callbacks {
			go callback(fields)
		}
	}
}
