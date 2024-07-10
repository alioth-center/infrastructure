package logger

type Option func(*customLogger)

func WithJsonFormatOpts() Option {
	return func(c *customLogger) {
		c.marshaller = defaultMarshaller
	}
}

func WithLevelOpts(level Level) Option {
	return func(c *customLogger) {
		c.level = level
	}
}

func WithStdWriterOpts() Option {
	return func(c *customLogger) {
		c.writer = NewStdoutConsoleWriter()
	}
}

func WithFileWriterOpts(file string) Option {
	return func(c *customLogger) {
		c.writer = NewFileWriter(file)
	}
}

func WithCustomWriterOpts(writer Writer) Option {
	return func(c *customLogger) {
		c.writer = writer
	}
}

func WithHookOpts(hook func(Fields), levels ...Level) Option {
	if len(levels) == 0 || hook == nil {
		return nil
	}

	return func(logger *customLogger) {
		for _, level := range levels {
			logger.hooks[level] = append(logger.hooks[level], hook)
		}
	}
}
