package logger

import (
	"context"
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
	//cw.Close()
}

func TestLogger(t *testing.T) {
	logger := Default()
	logger.Debug(NewFields(context.Background()).WithMessage("test1").WithData("hello"))
	logger.Info(NewFields(context.Background()).WithMessage("test2").WithData("hello"))
	logger.Warn(NewFields(context.Background()).WithMessage("test3").WithData("hello"))
	logger.Error(NewFields(context.Background()).WithMessage("test4").WithData("hello"))
	logger.Logf(LevelInfo, NewFields(context.Background()).WithMessage("test5"), "hello, %s", "world")
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
