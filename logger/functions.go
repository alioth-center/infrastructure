package logger

var defaultLogger = NewCustomLoggerWithOpts()

// Log logs a message at the specified level with the given fields.
//
// Parameters:
//
//	level (Level): The log level (e.g., Debug, Info, Warn, Error, Fatal, Panic).
//	fields (Fields): The fields to include in the log message.
func Log(level Level, fields Fields) {
	defaultLogger.Log(level, fields)
}

// Logf logs a formatted message at the specified level with the given fields.
//
// Parameters:
//
//	level (Level): The log level (e.g., Debug, Info, Warn, Error, Fatal, Panic).
//	fields (Fields): The fields to include in the log message.
//	format (string): The format string for the log message.
//	args (...any): The arguments for the format string.
func Logf(level Level, fields Fields, format string, args ...any) {
	defaultLogger.Logf(level, fields, format, args...)
}

// Debug logs a message at the Debug level with the given fields.
//
// Parameters:
//
//	fields (Fields): The fields to include in the log message.
func Debug(fields Fields) {
	defaultLogger.Debug(fields)
}

// Debugf logs a formatted message at the Debug level with the given fields.
//
// Parameters:
//
//	fields (Fields): The fields to include in the log message.
//	format (string): The format string for the log message.
//	args (...any): The arguments for the format string.
func Debugf(fields Fields, format string, args ...any) {
	defaultLogger.Debugf(fields, format, args...)
}

// Info logs a message at the Info level with the given fields.
//
// Parameters:
//
//	fields (Fields): The fields to include in the log message.
func Info(fields Fields) {
	defaultLogger.Info(fields)
}

// Infof logs a formatted message at the Info level with the given fields.
//
// Parameters:
//
//	fields (Fields): The fields to include in the log message.
//	format (string): The format string for the log message.
//	args (...any): The arguments for the format string.
func Infof(fields Fields, format string, args ...any) {
	defaultLogger.Infof(fields, format, args...)
}

// Warn logs a message at the Warn level with the given fields.
//
// Parameters:
//
//	fields (Fields): The fields to include in the log message.
func Warn(fields Fields) {
	defaultLogger.Warn(fields)
}

// Warnf logs a formatted message at the Warn level with the given fields.
//
// Parameters:
//
//	fields (Fields): The fields to include in the log message.
//	format (string): The format string for the log message.
//	args (...any): The arguments for the format string.
func Warnf(fields Fields, format string, args ...any) {
	defaultLogger.Warnf(fields, format, args...)
}

// Error logs a message at the Error level with the given fields.
//
// Parameters:
//
//	fields (Fields): The fields to include in the log message.
func Error(fields Fields) {
	defaultLogger.Error(fields)
}

// Errorf logs a formatted message at the Error level with the given fields.
//
// Parameters:
//
//	fields (Fields): The fields to include in the log message.
//	format (string): The format string for the log message.
//	args (...any): The arguments for the format string.
func Errorf(fields Fields, format string, args ...any) {
	defaultLogger.Errorf(fields, format, args...)
}

// Fatal logs a message at the Fatal level with the given fields and then exits the application.
//
// Parameters:
//
//	fields (Fields): The fields to include in the log message.
func Fatal(fields Fields) {
	defaultLogger.Fatal(fields)
}

// Fatalf logs a formatted message at the Fatal level with the given fields and then exits the application.
//
// Parameters:
//
//	fields (Fields): The fields to include in the log message.
//	format (string): The format string for the log message.
//	args (...any): The arguments for the format string.
func Fatalf(fields Fields, format string, args ...any) {
	defaultLogger.Fatalf(fields, format, args...)
}

// Panic logs a message at the Panic level with the given fields and then panics.
//
// Parameters:
//
//	fields (Fields): The fields to include in the log message.
func Panic(fields Fields) {
	defaultLogger.Panic(fields)
}

// Panicf logs a formatted message at the Panic level with the given fields and then panics.
//
// Parameters:
//
//	fields (Fields): The fields to include in the log message.
//	format (string): The format string for the log message.
//	args (...any): The arguments for the format string.
func Panicf(fields Fields, format string, args ...any) {
	defaultLogger.Panicf(fields, format, args...)
}
