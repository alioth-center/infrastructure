package http

import (
	"context"
	"github.com/alioth-center/infrastructure/utils/values"
	"net/http"
	"time"
)

type PreprocessedContext[request any, response any] interface {
	context.Context
	SetQueryParams(params Params)
	SetPathParams(params Params)
	SetHeaderParams(params Params)
	SetCookieParams(params Params)
	SetExtraParams(params Params)
	SetRawRequest(raw *http.Request)
	SetRequest(req request)
	SetRequestHeader(headers RequestHeader)
}

type Context[request any, response any] interface {
	context.Context

	// reset resets the acContext to its initial state.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		ctx.reset()
	// 	}
	reset()

	// setHandlers sets the handlers of the acContext.
	// example:
	//
	//	func initContext() {
	// 		chain := NewChain[request, response](
	//			func(ctx Context[request, response]) {
	//				// do something
	//				ctx.Next()
	//			},
	//			func(ctx Context[request, response]) {
	//				// do something
	//				ctx.Next()
	//			},
	//		)
	//		ctx.setHandlers(chain)
	//	}
	setHandlers(chain Chain[request, response])

	// Next calls the following handlers until abort or the end of the chain.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		ctx.Next()
	//		// do something after returned from the following handlers
	//	}
	Next()

	// Abort aborts the following handlers.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		ctx.Abort()
	//		// do something after aborted
	//	}
	Abort()

	// IsAborted returns true if the following handlers are aborted.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		if ctx.IsAborted() {
	//			// do something if aborted
	//		}
	//		// do something after aborted
	//	}
	IsAborted() bool

	// RawRequest returns the raw request.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		raw := ctx.RawRequest()
	//		// do something with raw
	//	}
	RawRequest() *http.Request

	// QueryParams returns the query params.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		params := ctx.QueryParams()
	//		// do something with params
	//	}
	QueryParams() Params

	// PathParams returns the path params.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		params := ctx.PathParams()
	//		// do something with params
	//	}
	PathParams() Params

	// HeaderParams returns the header params.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		params := ctx.HeaderParams()
	//		// do something with params
	//	}
	HeaderParams() Params

	// NormalHeaders returns the normal headers.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		headers := ctx.NormalHeaders()
	//		// do something with headers
	//	}
	NormalHeaders() RequestHeader

	// ExtraParams returns the extra params.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		params := ctx.ExtraParams()
	//		// do something with params
	//	}
	ExtraParams() Params

	// SetExtraParam sets the extra param.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		ctx.SetExtraParam("key", "value")
	//		// do something with params
	//	}
	SetExtraParam(key string, value string)

	// Request returns the processed request.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		req := ctx.Request()
	//		// do something with req
	//	}
	Request() request

	// Response returns the response.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		resp := ctx.Response()
	//		// do something with resp
	//	}
	Response() response

	// SetResponse sets the response.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		ctx.SetResponse(resp)
	//		// do something with resp
	//	}
	SetResponse(resp response)

	// ResponseHeaders returns the response headers.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		headers := ctx.ResponseHeaders()
	//		// do something with headers
	//	}
	ResponseHeaders() Params

	// SetResponseHeader sets the response header.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		ctx.SetResponseHeader("key", "value")
	//		// do something with resp
	//	}
	SetResponseHeader(key string, value string)

	// ResponseSetCookies returns the response cookies.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		cookies := ctx.ResponseSetCookies()
	//		// do something with cookies
	//	}
	ResponseSetCookies() []Cookie

	// SetResponseSetCookie sets the response cookie.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		ctx.SetResponseSetCookie(cookie)
	//		// do something with resp
	//	}
	SetResponseSetCookie(cookie Cookie)

	// StatusCode returns the status code.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		status := ctx.StatusCode()
	//		// do something with status
	//	}
	StatusCode() int

	// SetStatusCode sets the status code.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		ctx.SetStatusCode(StatusOK)
	//		// do something with status
	//	}
	SetStatusCode(status int)

	// Error returns the error.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		err := ctx.Error()
	//		// do something with err
	//	}
	Error() error

	// SetError sets the error.
	// example:
	//
	//	func handler(ctx Context[request, response]) {
	//		// do something
	//		ctx.SetError(err)
	//		// do something with err
	//	}
	SetError(err error)
}

type acContext[request any, response any] struct {
	queryParams  Params
	pathParams   Params
	headerParams Params
	cookieParams Params
	extraParams  Params

	idx int
	h   Chain[request, response]
	raw *http.Request
	ctx context.Context

	req        request
	resp       response
	setHeaders Params
	setCookies []Cookie
	headers    RequestHeader
	status     int
	err        error
}

func (c *acContext[request, response]) reset() {
	c.idx = -1
	c.h = []Handler[request, response]{}
}

func (c *acContext[request, response]) setHandlers(chain Chain[request, response]) {
	c.h = chain
}

func (c *acContext[request, response]) SetQueryParams(params Params) {
	c.queryParams = params
}

func (c *acContext[request, response]) SetPathParams(params Params) {
	c.pathParams = params
}

func (c *acContext[request, response]) SetHeaderParams(params Params) {
	c.headerParams = params
}

func (c *acContext[request, response]) SetCookieParams(params Params) {
	c.cookieParams = params
}

func (c *acContext[request, response]) SetExtraParams(params Params) {
	c.extraParams = params
}

func (c *acContext[request, response]) SetRawRequest(raw *http.Request) {
	c.raw = raw
}

func (c *acContext[request, response]) SetRequest(req request) {
	c.req = req
}

func (c *acContext[request, response]) SetRequestHeader(headers RequestHeader) {
	c.headers = headers
}

func (c *acContext[request, response]) Next() {
	c.idx++
	for c.idx < len(c.h) {
		c.h[c.idx](c)
		c.idx++
	}
}

func (c *acContext[request, response]) Abort() {
	c.idx = len(c.h)
}

func (c *acContext[request, response]) IsAborted() bool {
	return c.idx >= len(c.h)
}

func (c *acContext[request, response]) RawRequest() *http.Request {
	return c.raw
}

func (c *acContext[request, response]) QueryParams() Params {
	if c.queryParams == nil {
		c.queryParams = Params{}
	}

	return c.queryParams
}

func (c *acContext[request, response]) PathParams() Params {
	if c.pathParams == nil {
		c.pathParams = Params{}
	}

	return c.pathParams
}

func (c *acContext[request, response]) HeaderParams() Params {
	if c.headerParams == nil {
		c.headerParams = Params{}
	}

	return c.headerParams
}

func (c *acContext[request, response]) NormalHeaders() RequestHeader {
	return c.headers
}

func (c *acContext[request, response]) ExtraParams() Params {
	if c.extraParams == nil {
		c.extraParams = Params{}
	}

	return c.extraParams
}

func (c *acContext[request, response]) SetExtraParam(key string, value string) {
	if c.extraParams == nil {
		c.extraParams = Params{}
	}

	c.extraParams[key] = value
}

func (c *acContext[request, response]) Request() request {
	return c.req
}

func (c *acContext[request, response]) Response() response {
	return c.resp
}

func (c *acContext[request, response]) SetResponse(resp response) {
	c.resp = resp
}

func (c *acContext[request, response]) ResponseHeaders() Params {
	return c.setHeaders
}

func (c *acContext[request, response]) SetResponseHeader(key string, value string) {
	if c.setHeaders == nil {
		c.setHeaders = Params{}
	}

	c.setHeaders[key] = value
}

func (c *acContext[request, response]) ResponseSetCookies() []Cookie {
	return c.setCookies
}

func (c *acContext[request, response]) SetResponseSetCookie(cookie Cookie) {
	if c.setCookies == nil {
		c.setCookies = []Cookie{}
	}

	c.setCookies = append(c.setCookies, cookie)
}

func (c *acContext[request, response]) StatusCode() int {
	return c.status
}

func (c *acContext[request, response]) SetStatusCode(status int) {
	c.status = status
}

func (c *acContext[request, response]) Error() error {
	return c.err
}

func (c *acContext[request, response]) SetError(err error) {
	c.err = err
}

func (c *acContext[request, response]) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *acContext[request, response]) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *acContext[request, response]) Err() error {
	return c.ctx.Err()
}

func (c *acContext[request, response]) Value(key any) any {
	return c.ctx.Value(key)
}

type ContextOpts[request any, response any] func(ctx *acContext[request, response])

func WithContext[request any, response any](rawCtx context.Context) ContextOpts[request, response] {
	return func(ctx *acContext[request, response]) {
		ctx.ctx = rawCtx
	}
}

func WithQueryParams[request any, response any](params Params) ContextOpts[request, response] {
	return func(ctx *acContext[request, response]) {
		if params == nil {
			params = Params{}
		}

		ctx.queryParams = params
	}
}

func WithPathParams[request any, response any](params Params) ContextOpts[request, response] {
	return func(ctx *acContext[request, response]) {
		if params == nil {
			params = Params{}
		}

		ctx.pathParams = params
	}
}

func WithHeaderParams[request any, response any](params Params) ContextOpts[request, response] {
	return func(ctx *acContext[request, response]) {
		if params == nil {
			params = Params{}
		}

		ctx.headerParams = params
	}
}

func WithExtraParams[request any, response any](params Params) ContextOpts[request, response] {
	return func(ctx *acContext[request, response]) {
		if params == nil {
			params = Params{}
		}

		ctx.extraParams = params
	}
}

func WithCookieParams[request any, response any](params Params) ContextOpts[request, response] {
	return func(ctx *acContext[request, response]) {
		if params == nil {
			params = Params{}
		}

		ctx.cookieParams = params
	}
}

func WithRawRequest[request any, response any](raw *http.Request) ContextOpts[request, response] {
	return func(ctx *acContext[request, response]) {
		ctx.raw = raw
	}
}

func WithRequest[request any, response any](req request) ContextOpts[request, response] {
	return func(ctx *acContext[request, response]) {
		ctx.req = req
	}
}

func WithRequestHeader[request any, response any](headers RequestHeader) ContextOpts[request, response] {
	return func(ctx *acContext[request, response]) {
		ctx.headers = headers
	}
}

func WithResponse[request any, response any](resp response) ContextOpts[request, response] {
	return func(ctx *acContext[request, response]) {
		ctx.resp = resp
	}
}

func WithResponseHeaders[request any, response any](headers Params) ContextOpts[request, response] {
	return func(ctx *acContext[request, response]) {
		if headers == nil {
			headers = Params{}
		}

		ctx.setHeaders = headers
	}
}

func WithSetHeaders[request any, response any](headers Params) ContextOpts[request, response] {
	return func(ctx *acContext[request, response]) {
		if headers == nil {
			headers = Params{}
		}

		ctx.setHeaders = headers
	}
}

func WithSetCookies[request any, response any](cookies []Cookie) ContextOpts[request, response] {
	return func(ctx *acContext[request, response]) {
		if cookies == nil {
			cookies = []Cookie{}
		}

		ctx.setCookies = cookies
	}
}

func WithStatusCode[request any, response any](status int) ContextOpts[request, response] {
	return func(ctx *acContext[request, response]) {
		ctx.status = status
	}
}

func WithError[request any, response any](err error) ContextOpts[request, response] {
	return func(ctx *acContext[request, response]) {
		ctx.err = err
	}
}

func NewContext[request any, response any](opts ...ContextOpts[request, response]) Context[request, response] {
	ctx := &acContext[request, response]{
		queryParams:  Params{},
		pathParams:   Params{},
		headerParams: Params{},
		cookieParams: Params{},
		extraParams:  Params{},
		idx:          -1,
		h:            Chain[request, response]{},
		raw:          nil,
		ctx:          context.Background(),
		req:          values.Nil[request](),
		resp:         values.Nil[response](),
		setHeaders:   Params{},
		setCookies:   []Cookie{},
		headers:      RequestHeader{},
		status:       StatusOK,
		err:          nil,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(ctx)
		}
	}

	return ctx
}
