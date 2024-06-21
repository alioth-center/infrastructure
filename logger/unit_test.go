package logger

import (
	"context"
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	type data struct {
		Name string `json:"name" yaml:"name" xml:"name"`
		Age  int    `json:"age" yaml:"age" xml:"age"`
	}

	base := NewFields(context.Background()).WithMessage("test").WithLevel(LevelInfo).WithData(data{
		Name: "alice",
		Age:  114,
	})

	t.Run("JsonLog", func(t *testing.T) {
		f := NewFields(context.Background()).WithBaseFields(base)
		bytesOfLog, _ := marshalEntry(f.Export(), JsonMarshaller)
		t.Log(string(bytesOfLog))
	})

	t.Run("TextLog", func(t *testing.T) {
		f := NewFields(context.Background()).WithBaseFields(base)
		bytesOfLog, _ := marshalEntry(f.Export(), TextMarshaller)
		t.Log(string(bytesOfLog))
	})

	t.Run("CsvLog", func(t *testing.T) {
		f := NewFields(context.Background()).WithBaseFields(base)
		bytesOfLog, _ := marshalEntry(f.Export(), CsvMarshaller)
		t.Log(string(bytesOfLog))
	})

	t.Run("TsvLog", func(t *testing.T) {
		f := NewFields(context.Background()).WithBaseFields(base)
		bytesOfLog, _ := marshalEntry(f.Export(), TsvMarshaller)
		t.Log(string(bytesOfLog))
	})
}

func TestConsoleWriter(t *testing.T) {
	cw := ConsoleWriter()

	type data struct {
		Name string `json:"name" yaml:"name" xml:"name"`
		Age  int    `json:"age" yaml:"age" xml:"age"`
	}

	base := NewFields(context.Background()).WithMessage("test").WithLevel(LevelInfo).WithData(data{
		Name: "alice",
		Age:  114,
	})

	wb, _ := marshalEntry(NewFields(context.Background()).WithBaseFields(base).Export(), JsonMarshaller)
	cw.Write(wb)
	cw.Write(wb)
	cw.Write(wb)
	cw.Write(wb)
	// cw.Close()
}

func TestLogger(t *testing.T) {
	loggerTesting := func(logging Logger, t *testing.T) {
		logging.Debug(NewFields(context.Background()).WithMessage("test1").WithData("hello"))
		logging.Info(NewFields(context.Background()).WithMessage("test2").WithData("hello"))
		logging.Warn(NewFields(context.Background()).WithMessage("test3").WithData("hello"))
		logging.Error(NewFields(context.Background()).WithMessage("test4").WithData("hello"))
		logging.Debugf(NewFields(context.Background()).WithMessage("test5"), "hello, %s", "world")
		logging.Infof(NewFields(context.Background()).WithMessage("test6"), "hello, %s", "world")
		logging.Warnf(NewFields(context.Background()).WithMessage("test7"), "hello, %s", "world")
		logging.Errorf(NewFields(context.Background()).WithMessage("test8"), "hello, %s", "world")
		logging.Logf(LevelInfo, NewFields(context.Background()).WithMessage("test9"), "hello, %s", "world")
	}

	testingLoggers := []Logger{Default(), New(), Mute()}
	for i, ll := range testingLoggers {
		t.Run(fmt.Sprintf("%d:Loggeer", i), func(t *testing.T) {
			loggerTesting(ll, t)
		})
	}
}

func TestConfigConvert(t *testing.T) {
	cfg := Config{
		Level:     "info",
		Formatter: "json",
	}
	opt := convertConfigToOptions(cfg)
	logger := newLoggerWithOptions(opt)
	logger.Info(NewFields(context.Background()).WithMessage("test1").WithData("hello"))
}
