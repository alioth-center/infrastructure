package http

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/alioth-center/infrastructure/utils/network"

	"github.com/alioth-center/infrastructure/utils/values"
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

	// reset resets the acContext to its initial state, clearing handlers and resetting the index.
	reset()

	// setHandlers sets the handlers chain for the context.
	//
	// Parameters:
	//   chain (Chain[request, response]): The chain of handlers to be set.
	setHandlers(chain Chain[request, response])

	// Next calls the next handler in the chain.
	// If there are no more handlers, it stops execution.
	Next()

	// Abort stops the execution of the remaining handlers in the chain.
	Abort()

	// IsAborted checks if the execution of the remaining handlers has been aborted.
	//
	// Returns:
	//   bool: True if the execution has been aborted, otherwise false.
	IsAborted() bool

	// RawRequest returns the raw HTTP request.
	//
	// Returns:
	//   *http.Request: The raw HTTP request.
	RawRequest() *http.Request

	// QueryParams returns the query parameters.
	//
	// Returns:
	//   Params: The query parameters.
	QueryParams() Params

	// PathParams returns the path parameters.
	//
	// Returns:
	//   Params: The path parameters.
	PathParams() Params

	// HeaderParams returns the header parameters.
	//
	// Returns:
	//   Params: The header parameters.
	HeaderParams() Params

	// CookieParams returns the cookie parameters.
	//
	// Returns:
	//   Params: The cookie parameters.
	CookieParams() Params

	// NormalHeaders returns the normal request headers.
	// Includes following headers:
	//   - Accept
	//   - Accept-Encoding
	//   - Accept-Language
	//   - User-Agent
	//   - Content-Type
	//   - Content-Length
	//   - Origin
	//   - Referer
	//   - Authorization
	//   - ApiKey
	//
	// Returns:
	//   RequestHeader: The normal request headers.
	NormalHeaders() RequestHeader

	// ExtraParams returns the extra parameters.
	//
	// Returns:
	//   Params: The extra parameters.
	ExtraParams() Params

	// SetExtraParam sets a single extra parameter.
	//
	// Parameters:
	//   key (string): The key of the extra parameter.
	//   value (string): The value of the extra parameter.
	SetExtraParam(key string, value string)

	// Request returns the processed request.
	//
	// Returns:
	//   request: The processed request
	Request() request

	// Response returns the response.
	//
	// Returns:
	//   response: The response.
	Response() response

	// SetResponse sets the response.
	//
	// Parameters:
	//   resp (response): The response to be set.
	SetResponse(resp response)

	// ResponseHeaders returns the response headers.
	//
	// Returns:
	//   Params: The response headers.
	ResponseHeaders() Params

	// SetResponseHeader sets a single response header.
	//
	// Parameters:
	//   key (string): The key of the response header.
	//   value (string): The value of the response header.
	SetResponseHeader(key string, value string)

	// ResponseSetCookies returns the response set cookies.
	//
	// Returns:
	//   []Cookie: The response set cookies.
	ResponseSetCookies() []Cookie

	// SetResponseSetCookie sets a response cookie.
	//
	// Parameters:
	//   cookie (Cookie): The response cookie to be set.
	SetResponseSetCookie(cookie Cookie)

	// StatusCode returns the status code of the response.
	//
	// Returns:
	//   int: The status code of the response.
	StatusCode() int

	// SetStatusCode sets the status code of the response.
	//
	// Parameters:
	//   status (int): The status code to be set.
	SetStatusCode(status int)

	// ClientIP retrieves the client's IP address from the request context.
	// It checks various headers and the remote address in the following priority order:
	// 	1. X-Forwarded-For header (may contain multiple IPs, comma-separated)
	// 	2. X-Real-IP header
	// 	3. RemoteAddr from the raw request
	// 	4. Extra parameters passed in the context
	//
	// Returns:
	//   ip (string): The client's IP address if found and valid, otherwise an empty string.
	ClientIP() string

	// Error returns the error set in the context.
	//
	// Returns:
	//   error: The error set in the context.
	Error() error

	// SetError sets an error in the context.
	//
	// Parameters:
	//   err (error): The error to be set.
	SetError(err error)

	// SetValue sets a value in the context, equivalent to context.WithValue.
	//
	// Parameters:
	//   key (any): The key for the value.
	//   value (any): The value to be set.
	SetValue(key, value any)
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

func (c *acContext[request, response]) CookieParams() Params {
	if c.cookieParams == nil {
		c.cookieParams = Params{}
	}

	return c.cookieParams
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

func (c *acContext[request, response]) ClientIP() string {
	// check client IP from X-Forwarded-For
	if xForwardedFor := c.raw.Header.Get(HeaderXForwardedFor); xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		for _, ip := range ips {
			ip = strings.TrimSpace(ip)
			if network.IsValidIP(ip) {
				return ip
			}
		}
	}

	// check client IP from X-Real-IP
	if xRealIP := c.raw.Header.Get(HeaderXRealIP); xRealIP != "" && network.IsValidIP(xRealIP) {
		return xRealIP
	}

	// check client IP from RemoteAddr
	if remoteAddr := c.raw.RemoteAddr; remoteAddr != "" {
		if ip, _, err := net.SplitHostPort(remoteAddr); err == nil {
			if network.IsValidIP(ip) {
				return ip
			}
		}
	}

	// gin default client IP
	if clientIP, got := c.extraParams[RemoteIPKey]; got && network.IsValidIP(clientIP) {
		return clientIP
	}

	return ""
}

func (c *acContext[request, response]) Error() error {
	return c.err
}

func (c *acContext[request, response]) SetError(err error) {
	c.err = err
}

func (c *acContext[request, response]) SetValue(key, value any) {
	c.ctx = context.WithValue(c.ctx, key, value)
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
