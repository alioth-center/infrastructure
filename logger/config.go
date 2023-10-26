package logger

import (
	"strings"
)

type Config struct {
	Level          string `yaml:"level,omitempty" json:"level,omitempty" xml:"level,omitempty"`
	Formatter      string `yaml:"formatter,omitempty" json:"formatter,omitempty" xml:"formatter,omitempty"`
	StdoutFilePath string `yaml:"stdout_path,omitempty" json:"stdoutFilePath,omitempty" xml:"stdoutFilePath,omitempty"`
	StderrFilePath string `yaml:"stderr_path,omitempty" json:"stderrFilePath,omitempty" xml:"stderrFilePath,omitempty"`
}

func convertConfigToOptions(cfg Config) (opt Options) {
	opts := Options{
		StdoutWriter: ConsoleWriter(),
		StderrWriter: ConsoleErrorWriter(),
	}

	switch strings.ToLower(cfg.Level) {
	case "debug":
		opts.LogLevel = LevelDebug
	case "info":
		opts.LogLevel = LevelInfo
	case "warn":
		opts.LogLevel = LevelWarn
	case "error":
		opts.LogLevel = LevelError
	case "fatal":
		opts.LogLevel = LevelFatal
	case "panic":
		opts.LogLevel = LevelPanic
	default:
		opts.LogLevel = LevelInfo
	}

	switch strings.ToLower(cfg.Formatter) {
	case "json":
		opts.Marshaller = JsonMarshaller
	case "text":
		opts.Marshaller = TextMarshaller
	case "csv":
		opts.Marshaller = CsvMarshaller
	case "tsv":
		opts.Marshaller = TsvMarshaller
	default:
		opts.Marshaller = JsonMarshaller
	}

	if cfg.StdoutFilePath != "" {
		stdout, fwe := FileWriter(cfg.StdoutFilePath)
		if fwe == nil {
			opts.StdoutWriter = stdout
		}
	}

	if cfg.StderrFilePath != "" {
		stderr, fwe := FileWriter(cfg.StderrFilePath)
		if fwe == nil {
			opts.StderrWriter = stderr
		}
	}

	return opts
}
