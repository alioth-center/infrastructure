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

// EchoHandler return an echo handler.
// example:
//
//	chain := NewChain[request, request](
//		EchoHandler[request](),
//	)
func EchoHandler[payload any]() Handler[payload, payload] {
	return func(ctx Context[payload, payload]) {
		for k, v := range ctx.HeaderParams() {
			ctx.SetResponseHeader(k, v)
		}
		for k, v := range ctx.CookieParams() {
			ctx.SetResponseSetCookie(NewBasicCookie(k, v))
		}
		ctx.SetResponse(ctx.Request())
	}
}
