package logger

import (
	"fmt"
	"github.com/alioth-center/infrastructure/exit"
)

var (
	defaultLogger *logger = nil
)

type Options struct {
	LogLevel     Level
	Marshaller   Marshaller
	StdoutWriter Writer
	StderrWriter Writer
}

type Logger interface {
	init(options Options)
	Debug(fields Fields)
	Info(fields Fields)
	Warn(fields Fields)
	Error(fields Fields)
	Fatal(fields Fields)
	Panic(fields Fields)
	Log(level Level, fields Fields)
	Logf(level Level, fields Fields, format string, args ...any)
	Debugf(fields Fields, format string, args ...any)
	Infof(fields Fields, format string, args ...any)
	Warnf(fields Fields, format string, args ...any)
	Errorf(fields Fields, format string, args ...any)
	Fatalf(fields Fields, format string, args ...any)
	Panicf(fields Fields, format string, args ...any)
}

type logger struct {
	options Options
}

func (l *logger) init(options Options) {
	l.options = Options{
		LogLevel:     options.LogLevel,
		Marshaller:   options.Marshaller,
		StdoutWriter: options.StdoutWriter,
		StderrWriter: options.StderrWriter,
	}

	if l.options.LogLevel == "" {
		l.options.LogLevel = LevelInfo
	}

	if l.options.Marshaller == nil {
		l.options.Marshaller = JsonMarshaller
	}

	if l.options.StdoutWriter == nil {
		l.options.StdoutWriter = ConsoleWriter()
	}

	if l.options.StderrWriter == nil {
		l.options.StderrWriter = ConsoleErrorWriter()
	}
}

func (l *logger) marshalFieldsToBytes(fields Fields) (data []byte) {
	d, _ := marshalEntry(fields.export(), l.options.Marshaller)
	return d
}

func (l *logger) Log(level Level, fields Fields) {
	fields = fields.WithLevel(level)
	switch level {
	case LevelDebug:
		l.Debug(fields)
	case LevelInfo:
		l.Info(fields)
	case LevelWarn:
		l.Warn(fields)
	case LevelError:
		l.Error(fields)
	case LevelFatal:
		l.Fatal(fields)
	case LevelPanic:
		l.Panic(fields)
	default:
		l.Info(fields)
	}
}

func (l *logger) Logf(level Level, fields Fields, format string, args ...any) {
	l.Log(level, fields.WithMessage(fmt.Sprintf(format, args...)))
}

func (l *logger) Debug(fields Fields) {
	if levelValueMap[l.options.LogLevel] <= levelValueMap[LevelDebug] {
		l.options.StdoutWriter.Write(l.marshalFieldsToBytes(fields.WithLevel(LevelDebug)))
	}
}

func (l *logger) Info(fields Fields) {
	if levelValueMap[l.options.LogLevel] <= levelValueMap[LevelInfo] {
		l.options.StdoutWriter.Write(l.marshalFieldsToBytes(fields.WithLevel(LevelInfo)))
	}
}

func (l *logger) Warn(fields Fields) {
	if levelValueMap[l.options.LogLevel] <= levelValueMap[LevelWarn] {
		l.options.StdoutWriter.Write(l.marshalFieldsToBytes(fields.WithLevel(LevelWarn)))
		l.options.StderrWriter.Write(l.marshalFieldsToBytes(fields.WithLevel(LevelWarn)))
	}
}

func (l *logger) Error(fields Fields) {
	if levelValueMap[l.options.LogLevel] <= levelValueMap[LevelError] {
		l.options.StdoutWriter.Write(l.marshalFieldsToBytes(fields.WithLevel(LevelError)))
		l.options.StderrWriter.Write(l.marshalFieldsToBytes(fields.WithLevel(LevelError)))
	}
}

func (l *logger) Fatal(fields Fields) {
	if levelValueMap[l.options.LogLevel] <= levelValueMap[LevelFatal] {
		l.options.StdoutWriter.Write(l.marshalFieldsToBytes(fields.WithLevel(LevelFatal)))
		l.options.StderrWriter.Write(l.marshalFieldsToBytes(fields.WithLevel(LevelFatal)))
	}
}

func (l *logger) Panic(fields Fields) {
	if levelValueMap[l.options.LogLevel] <= levelValueMap[LevelPanic] {
		l.options.StdoutWriter.Write(l.marshalFieldsToBytes(fields.WithLevel(LevelPanic)))
		l.options.StderrWriter.Write(l.marshalFieldsToBytes(fields.WithLevel(LevelPanic)))
	}
}

func (l *logger) Debugf(fields Fields, format string, args ...any) {
	l.Debug(fields.WithMessage(fmt.Sprintf(format, args...)))
}

func (l *logger) Infof(fields Fields, format string, args ...any) {
	l.Info(fields.WithMessage(fmt.Sprintf(format, args...)))
}

func (l *logger) Warnf(fields Fields, format string, args ...any) {
	l.Warn(fields.WithMessage(fmt.Sprintf(format, args...)))
}

func (l *logger) Errorf(fields Fields, format string, args ...any) {
	l.Error(fields.WithMessage(fmt.Sprintf(format, args...)))
}

func (l *logger) Fatalf(fields Fields, format string, args ...any) {
	l.Fatal(fields.WithMessage(fmt.Sprintf(format, args...)))
}

func (l *logger) Panicf(fields Fields, format string, args ...any) {
	l.Panic(fields.WithMessage(fmt.Sprintf(format, args...)))
}

func New() Logger {
	var l = &logger{}
	l.init(Options{
		LogLevel:     LevelInfo,
		Marshaller:   JsonMarshaller,
		StdoutWriter: ConsoleWriter(),
		StderrWriter: ConsoleErrorWriter(),
	})
	return l
}

func Default() Logger {
	if defaultLogger == nil {
		defaultLogger = &logger{}
		defaultLogger.init(Options{
			LogLevel:     LevelInfo,
			Marshaller:   JsonMarshaller,
			StdoutWriter: ConsoleWriter(),
			StderrWriter: ConsoleErrorWriter(),
		})
	}
	return defaultLogger
}

func NewLoggerWithOptions(options Options) Logger {
	var l = &logger{}
	exit.Register(func(_ string) string {
		options.StderrWriter.Close()
		options.StdoutWriter.Close()
		return "logger stopped"
	}, "logger")
	l.init(options)
	return l
}

func NewLoggerWithConfig(cfg Config) Logger {
	return NewLoggerWithOptions(convertConfigToOptions(cfg))
}
