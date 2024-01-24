package http

type Chain[request any, response any] []Handler[request, response]

// AddHandlerFront adds handlers to the front of the handler chain.
// example:
//
//	chain.AddHandlerFront(
//		func(ctx Context[request, response]) {
//			// do something
//			ctx.Next()
//		},
//		func(ctx Context[request, response]) {
//			// do something
//			ctx.Next()
//		},
//	)
func (c Chain[request, response]) AddHandlerFront(h ...Handler[request, response]) Chain[request, response] {
	return append(h, c...)
}

// AddHandlerBack adds handlers to the back of the handler chain.
// example:
//
//	chain.AddHandlerBack(
//		func(ctx Context[request, response]) {
//			// do something
//			ctx.Next()
//		},
//		func(ctx Context[request, response]) {
//			// do something
//			ctx.Next()
//		},
//	)
func (c Chain[request, response]) AddHandlerBack(h ...Handler[request, response]) Chain[request, response] {
	return append(c, h...)
}

// Execute executes the handler chain.
func (c Chain[request, response]) Execute(ctx Context[request, response]) {
	ctx.reset()
	ctx.setHandlers(c)
	ctx.Next()
}

// NewChain creates a new handler chain.
// example:
//
//	chain := NewChain[request, response](
//		func(ctx Context[request, response]) {
//			// do something
//			ctx.Next()
//		},
//		func(ctx Context[request, response]) {
//			// do something
//			ctx.Next()
//		},
//	)
func NewChain[request any, response any](handlers ...Handler[request, response]) Chain[request, response] {
	if handlers == nil {
		return []Handler[request, response]{EmptyHandler[request, response]()}
	}

	chain := make(Chain[request, response], 0, len(handlers))
	for _, handler := range handlers {
		if handler != nil {
			chain = append(chain, handler)
		}
	}
	if len(chain) == 0 {
		return []Handler[request, response]{EmptyHandler[request, response]()}
	}

	return chain
}
