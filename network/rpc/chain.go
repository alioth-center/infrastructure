package rpc

type Chain[request any, response any] []Handler[request, response]

func (c Chain[request, response]) AddHandlerFront(h ...Handler[request, response]) Chain[request, response] {
	return append(h, c...)
}

func (c Chain[request, response]) AddHandlerBack(h ...Handler[request, response]) Chain[request, response] {
	return append(c, h...)
}

func (c Chain[request, response]) Run(ctx *Context[request, response]) {
	ctx.hc = c
	ctx.idx = -1
	ctx.Next()
}

func NewChain[request any, response any](handlers ...Handler[request, response]) Chain[request, response] {
	if handlers == nil {
		return []Handler[request, response]{}
	}

	chain := make(Chain[request, response], 0, len(handlers))
	for _, handler := range handlers {
		if handler != nil {
			chain = append(chain, handler)
		}
	}

	return chain
}
