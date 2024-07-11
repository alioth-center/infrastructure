package logger

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

func Default() Logger {
	return defaultLogger
}

func Mute() Logger {
	return &customLogger{level: LevelInfo, hooks: map[Level][]func(Fields){}}
}

func File(file string, level Level) Logger {
	return NewCustomLoggerWithOpts(WithFileWriterOpts(file), WithLevelOpts(level), WithJsonFormatOpts())
}
