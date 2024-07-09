package rpc

import (
	"fmt"
	"net"
	"os"

	"github.com/alioth-center/infrastructure/exit"
	"google.golang.org/grpc"
)

type Engine struct {
	serving    bool
	connection net.Listener
	server     *grpc.Server
	services   []Service
	address    string
}

func (e *Engine) Serving() bool {
	return e.serving
}

func (e *Engine) Serve(address string) (err error) {
	if e.serving {
		return NewServerAlreadyServingError(e.address)
	}

	e.serving = true
	e.address = address

	// 只有在未启动服务时才需要注册退出
	exit.RegisterExitEvent(func(_ os.Signal) {
		e.server.GracefulStop()
		fmt.Println("rpc engine stopped")
	}, "SHUTDOWN_RPC_SERVER")

	conn, listenErr := net.Listen("tcp", address)
	defer func() {
		_ = conn.Close()
	}()

	if listenErr != nil {
		return fmt.Errorf("rpc server failed to listen: %w", listenErr)
	}
	e.connection = conn
	e.server = grpc.NewServer()

	for _, service := range e.services {
		service.Initialization()
		service.BindEngine(e.server)
	}

	return e.server.Serve(e.connection)
}

func (e *Engine) ServeAsync(address string, ex chan struct{}) (exitChan chan error) {
	ec := make(chan error)
	if e.serving {
		ec <- NewServerAlreadyServingError(e.address)
		return ec
	}

	// 只有在未启动服务时才需要退出
	exit.RegisterExitEvent(func(signal os.Signal) {
		ex <- struct{}{}
		fmt.Println("rpc server stopped")
	}, "SHUTDOWN_RPC_SERVER")

	go func() {
		select {
		case ec <- e.Serve(address):
			return
		case <-ex:
			return
		}
	}()

	return ec
}

func (e *Engine) AddService(services ...Service) {
	e.services = append(e.services, services...)
}

func NewEngine() *Engine {
	return &Engine{
		serving:    false,
		services:   make([]Service, 0),
		server:     nil,
		connection: nil,
		address:    "",
	}
}
