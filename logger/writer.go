package logger

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/alioth-center/infrastructure/exit"
	"github.com/alioth-center/infrastructure/utils/concurrency"
)

type Writer interface {
	Write(data []byte)
	Close()
}

var fileWriters = concurrency.NewMap[string, Writer]()

type fileLogWriter struct {
	f      *os.File
	buffer chan []byte
	closed atomic.Bool
}

func (fw *fileLogWriter) Write(data []byte) {
	if !fw.closed.Load() {
		fw.buffer <- data
	}
}

func (fw *fileLogWriter) Close() {
	if !fw.closed.Load() {
		fw.closed.Store(true)
		close(fw.buffer)
		for data := range fw.buffer {
			_, _ = fw.f.Write(data)
		}
		_ = fw.f.Close()
	}
}

func (fw *fileLogWriter) serve() {
	exit.RegisterExitEvent(func(_ os.Signal) {
		fw.Close()
	}, "EXIT_FILE_LOGGER:"+fw.f.Name())

	for data := range fw.buffer {
		_, _ = fw.f.Write(data)
	}
}

func NewFileWriter(path string) Writer {
	// exist file writer, return it
	if w, ok := fileWriters.Get(path); ok {
		return w
	}

	f, e := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if e != nil {
		return nil
	}

	w := &fileLogWriter{
		f:      f,
		buffer: make(chan []byte, 1024),
		closed: atomic.Bool{},
	}
	w.closed.Store(false)
	go w.serve()

	fileWriters.Set(path, w)

	return w
}

type rotationFileWriter struct {
	baseDir  string
	rotation func(time.Time) string
	lastFile atomic.Value
}

func (r *rotationFileWriter) Write(data []byte) {
	logFile := filepath.Join(r.baseDir, r.rotation(time.Now()))
	if r.lastFile.Load() != logFile {
		lastWriter, exist := fileWriters.Get(logFile)
		if exist && lastWriter != nil {
			lastWriter.Close()
		}
	}

	writer := NewFileWriter(logFile)
	writer.Write(data)
}

func (r *rotationFileWriter) Close() {
	lastFile := r.lastFile.Load().(string)
	lastWriter, exist := fileWriters.Get(lastFile)
	if exist && lastWriter != nil {
		lastWriter.Close()
	}
}

func NewTimeBasedRotationFileWriter(directory string, rotation func(time time.Time) (filename string)) Writer {
	atValue := atomic.Value{}
	atValue.Store(filepath.Join(directory, rotation(time.Now())))
	return &rotationFileWriter{
		baseDir:  directory,
		rotation: rotation,
		lastFile: atValue,
	}
}

type consoleWriter struct {
	console *os.File
}

func (c consoleWriter) Write(data []byte) {
	if c.console == os.Stdout || c.console == os.Stderr {
		_, _ = c.console.Write(data)
	}
}

func (c consoleWriter) Close() {}

func NewStdoutConsoleWriter() Writer {
	return consoleWriter{
		console: os.Stdout,
	}
}

func NewStderrConsoleWriter() Writer {
	return consoleWriter{
		console: os.Stderr,
	}
}

type multiWriter struct {
	writers []Writer
}

func (m multiWriter) Write(data []byte) {
	for _, writer := range m.writers {
		writer.Write(data)
	}
}

func (m multiWriter) Close() {
	for _, writer := range m.writers {
		writer.Close()
	}
}

func NewMultiWriter(writers ...Writer) Writer {
	return multiWriter{
		writers: writers,
	}
}
