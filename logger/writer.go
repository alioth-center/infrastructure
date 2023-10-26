package logger

import (
	"fmt"
	"github.com/alioth-center/infrastructure/errors"
	"os"
	"sync"
)

var (
	consoleWriter    Writer
	consoleErrWriter Writer
)

type Writer interface {
	init(output *os.File)
	errors() chan error
	Write(data []byte)
	Close()
}

type writer struct {
	iw     *os.File
	ex     chan struct{}
	es     chan error
	closed bool
	buffer chan []byte
	wg     *sync.WaitGroup
}

func (w *writer) init(output *os.File) {
	if output == nil {
		output = os.Stdout
	}

	w.iw = output
	w.ex = make(chan struct{})
	w.es = make(chan error, 16)
	w.closed = false
	w.buffer = make(chan []byte, 1024)
	w.wg = &sync.WaitGroup{}

	go w.serve()
}

func (w *writer) serve() {
	defer w.wg.Done()
	for {
		select {
		case <-w.ex:
			w.closed = true

			// 等待缓冲区中的数据全部写入
			for {
				select {
				case data := <-w.buffer:
					_, e := w.iw.Write(data)
					if e != nil {
						w.es <- e
					}
				default:
					goto end
				}
			}
		end:
			close(w.buffer)
			close(w.es)
			return
		case data := <-w.buffer:
			_, e := w.iw.Write(data)
			if e != nil {
				w.es <- e
			}
		}
	}
}

func (w *writer) errors() chan error {
	return w.es
}

func (w *writer) Write(data []byte) {
	if w.closed {
		return
	}

	w.buffer <- data
}

func (w *writer) Close() {
	if w.closed {
		return
	}

	w.wg.Add(1)
	w.ex <- struct{}{}
	_ = w.iw.Close()

	// 等待缓冲区中的数据写入
	w.wg.Wait()
}

func ConsoleWriter() Writer {
	return consoleWriter
}

func ConsoleErrorWriter() Writer {
	return consoleErrWriter
}

func FileWriter(path string) (w Writer, err error) {
	fi, fie := os.Stat(path)
	if fie != nil {
		if os.IsNotExist(fie) {
			// 文件不存在，创建文件
			f, cfe := os.Create(path)
			if cfe != nil {
				return nil, fmt.Errorf("create log file error: %w", cfe)
			}
			w = &writer{}
			w.init(f)
			return w, nil
		}

		return nil, fie
	}

	if fi.IsDir() {
		// 文件路径是一个目录，返回错误
		return nil, errors.NewFileWriterWriteToDirectoryError(path)
	}

	f, ope := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if ope != nil {
		// 打开文件失败，返回错误
		return nil, fmt.Errorf("open log file error: %w", ope)
	}

	w = &writer{}
	w.init(f)
	return w, nil
}
