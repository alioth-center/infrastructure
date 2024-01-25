package http

type Handler[request any, response any] func(ctx Context[request, response])

// EmptyHandler return an empty handler.
// example:
//
//	chain := NewChain[request, response](
//		EmptyHandler[request, response](),
//		func(ctx Context[request, response]) {
//			// do something
//			ctx.Next()
//		},
//	)
func EmptyHandler[request any, response any]() Handler[request, response] {
	return func(ctx Context[request, response]) {}
}
