package logger

import (
	"os"
)

func init() {
	consoleWriter = &writer{}
	consoleWriter.init(os.Stdout)
	consoleErrWriter = &writer{}
	consoleErrWriter.init(os.Stderr)

	defaultLogger = &logger{}
	defaultLogger.init(Options{
		LogLevel:     LevelInfo,
		Marshaller:   JsonMarshaller,
		StdoutWriter: ConsoleWriter(),
		StderrWriter: ConsoleErrorWriter(),
	})
}

func Log(level Level, fields Fields) {
	defaultLogger.Log(level, fields)
}

func Logf(level Level, fields Fields, format string, args ...any) {
	defaultLogger.Logf(level, fields, format, args...)
}

func Debug(fields Fields) {
	defaultLogger.Debug(fields)
}

func Debugf(fields Fields, format string, args ...any) {
	defaultLogger.Debugf(fields, format, args...)
}

func Info(fields Fields) {
	defaultLogger.Info(fields)
}

func Infof(fields Fields, format string, args ...any) {
	defaultLogger.Infof(fields, format, args...)
}

func Warn(fields Fields) {
	defaultLogger.Warn(fields)
}

func Warnf(fields Fields, format string, args ...any) {
	defaultLogger.Warnf(fields, format, args...)
}

func Error(fields Fields) {
	defaultLogger.Error(fields)
}

func Errorf(fields Fields, format string, args ...any) {
	defaultLogger.Errorf(fields, format, args...)
}

func Fatal(fields Fields) {
	defaultLogger.Fatal(fields)
}

func Fatalf(fields Fields, format string, args ...any) {
	defaultLogger.Fatalf(fields, format, args...)
}

func Panic(fields Fields) {
	defaultLogger.Panic(fields)
}

func Panicf(fields Fields, format string, args ...any) {
	defaultLogger.Panicf(fields, format, args...)
}

func SetLevel(level Level) {
	defaultLogger.options.LogLevel = level
}

func SetMarshaller(marshaller Marshaller) {
	defaultLogger.options.Marshaller = marshaller
}
