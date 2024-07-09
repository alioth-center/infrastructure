package logger

import (
	"fmt"
	"os"

	"github.com/alioth-center/infrastructure/exit"
)

var defaultLogger *logger = nil

type Options struct {
	LogLevel     Level
	Marshaller   Marshaller
	StdoutWriter Writer
	StderrWriter Writer
	notStdout    bool
	notStderr    bool
}

type Logger interface {
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
	d, _ := marshalEntry(fields.Export(), l.options.Marshaller)
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
	if LevelValueMap[l.options.LogLevel] <= LevelValueMap[LevelDebug] {
		l.options.StdoutWriter.Write(l.marshalFieldsToBytes(fields.WithLevel(LevelDebug)))
	}
}

func (l *logger) Info(fields Fields) {
	if LevelValueMap[l.options.LogLevel] <= LevelValueMap[LevelInfo] {
		l.options.StdoutWriter.Write(l.marshalFieldsToBytes(fields.WithLevel(LevelInfo)))
	}
}

func (l *logger) Warn(fields Fields) {
	if LevelValueMap[l.options.LogLevel] <= LevelValueMap[LevelWarn] {
		l.options.StdoutWriter.Write(l.marshalFieldsToBytes(fields.WithLevel(LevelWarn)))
		l.options.StderrWriter.Write(l.marshalFieldsToBytes(fields.WithLevel(LevelWarn)))
	}
}

func (l *logger) Error(fields Fields) {
	if LevelValueMap[l.options.LogLevel] <= LevelValueMap[LevelError] {
		l.options.StdoutWriter.Write(l.marshalFieldsToBytes(fields.WithLevel(LevelError)))
		l.options.StderrWriter.Write(l.marshalFieldsToBytes(fields.WithLevel(LevelError)))
	}
}

func (l *logger) Fatal(fields Fields) {
	if LevelValueMap[l.options.LogLevel] <= LevelValueMap[LevelFatal] {
		l.options.StdoutWriter.Write(l.marshalFieldsToBytes(fields.WithLevel(LevelFatal)))
		l.options.StderrWriter.Write(l.marshalFieldsToBytes(fields.WithLevel(LevelFatal)))
	}
}

func (l *logger) Panic(fields Fields) {
	if LevelValueMap[l.options.LogLevel] <= LevelValueMap[LevelPanic] {
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

type muteLogger struct{}

func (m muteLogger) Debug(fields Fields) {}

func (m muteLogger) Info(fields Fields) {}

func (m muteLogger) Warn(fields Fields) {}

func (m muteLogger) Error(fields Fields) {}

func (m muteLogger) Fatal(fields Fields) {}

func (m muteLogger) Panic(fields Fields) {}

func (m muteLogger) Log(level Level, fields Fields) {}

func (m muteLogger) Logf(level Level, fields Fields, format string, args ...any) {}

func (m muteLogger) Debugf(fields Fields, format string, args ...any) {}

func (m muteLogger) Infof(fields Fields, format string, args ...any) {}

func (m muteLogger) Warnf(fields Fields, format string, args ...any) {}

func (m muteLogger) Errorf(fields Fields, format string, args ...any) {}

func (m muteLogger) Fatalf(fields Fields, format string, args ...any) {}

func (m muteLogger) Panicf(fields Fields, format string, args ...any) {}

func New() Logger {
	l := &logger{}
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

func Mute() Logger {
	return &muteLogger{}
}

func newLoggerWithOptions(options Options) Logger {
	l := &logger{}
	exit.RegisterExitEvent(func(signal os.Signal) {
		// stdout 和 stderr 不能关闭，可能会导致异常情况
		if options.notStdout {
			options.StdoutWriter.Close()
		}
		if options.notStderr {
			options.StderrWriter.Close()
		}
		fmt.Println("logger stopped")
	}, "CLEANUP_LOGGER")
	l.init(options)
	return l
}

func NewLoggerWithConfig(cfg Config) Logger {
	return newLoggerWithOptions(convertConfigToOptions(cfg))
}

func NewLoggerWithCustomWriter(stdout, stderr Writer, closable bool, formatter Marshaller, level Level) Logger {
	l := &logger{}
	if closable {
		exit.RegisterExitEvent(func(signal os.Signal) {
			stdout.Close()
			stderr.Close()
			fmt.Println("logger stopped")
		}, "CLEANUP_LOGGER")
	}

	l.init(Options{
		LogLevel:     level,
		Marshaller:   formatter,
		StdoutWriter: stdout,
		StderrWriter: stderr,
	})

	return l
}
