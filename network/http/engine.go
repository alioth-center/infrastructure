package http

import (
	"fmt"
	"os"

	"github.com/alioth-center/infrastructure/exit"
	"github.com/alioth-center/infrastructure/trace"
	"github.com/alioth-center/infrastructure/utils/values"
	"github.com/gin-gonic/gin"
)

type Engine struct {
	core    *gin.Engine
	serving bool

	baseRouter  Router
	endpoints   []EndPointInterface
	middlewares []gin.HandlerFunc
}

func (e *Engine) registerEndpoints() {
	if len(e.middlewares) > 0 {
		e.core.Use(e.middlewares...)
	}

	for _, ep := range e.endpoints {
		ep.bindRouter(e.core.Group(""), e.baseRouter)
	}
}

func (e *Engine) BaseRouter() Router {
	return e.baseRouter
}

func (e *Engine) AddEndPoints(eps ...EndPointInterface) {
	if e.endpoints == nil {
		e.endpoints = []EndPointInterface{}
	}

	e.endpoints = append(e.endpoints, eps...)
}

func (e *Engine) AddMiddlewares(middleware ...gin.HandlerFunc) {
	if e.middlewares == nil {
		e.middlewares = []gin.HandlerFunc{}
	}

	e.middlewares = append(e.middlewares, middleware...)
}

func (e *Engine) Serve(bindAddress string) error {
	if e.serving {
		return ServerAlreadyServingError{Address: bindAddress}
	}

	e.registerEndpoints()
	return e.core.Run(bindAddress)
}

func (e *Engine) ServeAsync(bindAddress string, exitChan chan struct{}) (errChan chan error) {
	ec := make(chan error)
	if e.serving {
		ec <- ServerAlreadyServingError{Address: bindAddress}
		return ec
	}

	exit.RegisterExitEvent(func(_ os.Signal) {
		exitChan <- struct{}{}
		fmt.Println("http server stopped")
	}, "SHUTDOWN_HTTP_SERVER")

	go func() {
		select {
		case ec <- e.Serve(bindAddress):
			return
		case <-exitChan:
			return
		}
	}()

	return ec
}

func (e *Engine) traceContext(ctx *gin.Context) {
	tid := ctx.GetHeader(TraceHeaderKey())
	if tid == "" {
		tid = trace.GetTid(trace.NewContext())
	}
	ctx.Set(trace.ContextKey(), tid)
}

func (e *Engine) defaultHandler(ctx *gin.Context) {
	errResponse := &FrameworkResponse{
		ErrorCode:    ErrorCodeResourceNotFound,
		ErrorMessage: values.BuildStringsWithJoin(" ", "resource", ctx.Request.URL.String(), "not found or function not implemented"),
		RequestID:    trace.GetTid(ctx),
		Data:         nil,
	}

	ctx.JSON(StatusNotFound, errResponse)
}

func NewEngine(base string) *Engine {
	e := &Engine{
		core:        gin.New(),
		endpoints:   []EndPointInterface{},
		baseRouter:  NewRouter(base),
		middlewares: []gin.HandlerFunc{},
	}

	e.core.Use(gin.Recovery())
	e.core.Use(e.traceContext)
	e.core.NoRoute(e.defaultHandler)
	e.core.NoMethod(e.defaultHandler)

	return e
}
