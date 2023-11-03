package rpc

import "github.com/alioth-center/infrastructure/logger"

type Chain[request any, response any] []Handler[request, response]

func (c Chain[request, response]) AddHandler(h ...Handler[request, response]) Chain[request, response] {
	return append(c, h...)
}

func (c Chain[request, response]) Run(ctx *Context[request, response]) {
	ctx.h = c
	ctx.idx = -1
	ctx.Next()
}

func NewChain[request any, response any](handlers ...Handler[request, response]) Chain[request, response] {
	if handlers == nil {
		handlers = make([]Handler[request, response], 0)
	}
	return handlers
}

func DefaultChain[request any, response any](logger logger.Logger, serviceName string, handlerName string, handlers ...Handler[request, response]) Chain[request, response] {
	if handlers == nil {
		handlers = make([]Handler[request, response], 0)
	}
	return NewChain[request, response](
		RecoveryHandler[request, response](logger, serviceName, handlerName),
		LogInputAndOutputHandler[request, response](logger, serviceName, handlerName),
	).AddHandler(
		handlers...,
	)
}
