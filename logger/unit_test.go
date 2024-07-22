package logger

import (
	"github.com/alioth-center/infrastructure/trace"
	"testing"
	"time"
)

func TestCustom(t *testing.T) {
	logger := NewCustomLoggerWithOpts(WithCustomWriterOpts(NewTimeBasedRotationFileWriter("./", func(time time.Time) (filename string) { return time.Format("20060102_15") + ".jsonl" })), WithLevelOpts(LevelDebug), WithJsonFormatOpts())
	logger.Info(NewFields(trace.NewContext()).WithMessage("hello world"))
	time.Sleep(time.Second)
}

func TestLoggerFunction(t *testing.T) {
	ctx := trace.NewContext()
	base := NewFields(ctx).WithMessage("test").WithData(map[string]any{"foo": "bar", "nav": 0.1}).WithField("field", "value").WithCallTime(time.Now()).WithTraceID(trace.GetTid(ctx)).WithLevel(LevelInfo)
	Debug(NewFields(ctx).WithBaseFields(base))
	Debugf(NewFields(ctx).WithBaseFields(base), "test %s", "format")
	Info(NewFields(ctx).WithBaseFields(base))
	Infof(NewFields(ctx).WithBaseFields(base), "test %s", "format")
	Warn(NewFields(ctx).WithBaseFields(base))
	Warnf(NewFields(ctx).WithBaseFields(base), "test %s", "format")
	Error(NewFields(ctx).WithBaseFields(base))
	Errorf(NewFields(ctx).WithBaseFields(base), "test %s", "format")
	Fatal(NewFields(ctx).WithBaseFields(base))
	Fatalf(NewFields(ctx).WithBaseFields(base), "test %s", "format")
	Panic(NewFields(ctx).WithBaseFields(base))
	Panicf(NewFields(ctx).WithBaseFields(base), "test %s", "format")
	Log(LevelInfo, NewFields(ctx).WithBaseFields(base))
	Logf(LevelInfo, NewFields(ctx).WithBaseFields(base), "test %s", "format")
}

func TestLoggerCtor(t *testing.T) {
	logger := Mute()
	logger = Default()
	logger = File("./test.log", LevelInfo)
	logger.Info(NewFields().WithMessage("test"))
}

func TestWriter(t *testing.T) {
	writer := NewStdoutConsoleWriter()
	writer.Write([]byte("hello world"))
	writer.Close()

	writer = NewMultiWriter(NewStdoutConsoleWriter(), NewStderrConsoleWriter())
	writer.Write([]byte("hello world"))
	writer.Close()

	writer = NewFileWriter("./test_writer.log")
	writer.Write([]byte("hello world"))
	writer.Close()

	writer = NewTimeBasedRotationFileWriter("./", func(time time.Time) (_ string) { return "test_timed.jsonl" })
	writer.Write([]byte("hello world"))
	writer.Close()
}
