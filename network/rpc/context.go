package rpc

import (
	"context"
	"time"

	"github.com/alioth-center/infrastructure/trace"
)

const (
	abortIndex int = 1 << 31
)

type Context[request any, response any] struct {
	idx  int
	ctx  context.Context
	hc   Chain[request, response]
	err  error
	req  request
	resp response
}

func (c *Context[request, response]) Next() {
	c.idx++
	for c.idx < len(c.hc) {
		c.hc[c.idx](c)
		c.idx++
	}
}

func (c *Context[request, response]) Abort() {
	c.idx = abortIndex
}

func (c *Context[request, response]) IsAborted() bool {
	return c.idx >= abortIndex
}

func (c *Context[request, response]) GetRequest() request {
	return c.req
}

func (c *Context[request, response]) GetResponse() response {
	return c.resp
}

func (c *Context[request, response]) GetContext() context.Context {
	return c.ctx
}

func (c *Context[request, response]) SetRequest(req request) {
	c.req = req
}

func (c *Context[request, response]) SetResponse(resp response) {
	c.resp = resp
}

func (c *Context[request, response]) SetContext(ctx context.Context) {
	c.ctx = trace.FromContext(ctx)
}

func (c *Context[request, response]) Error() error {
	return c.err
}

func (c *Context[request, response]) SetError(err error) {
	c.err = err
}

func (c *Context[request, response]) SetResult(resp response, err error) {
	c.resp, c.err = resp, err
}

func (c *Context[request, response]) TraceID() string {
	return trace.GetTid(c.ctx)
}

func (c *Context[request, response]) GetContextClientIP() (ip string, err error) {
	ip = trace.GetClientIPFromPeer(c.ctx)
	if ip == "" {
		return "", NewGetRPCClientIPFailedError()
	}

	return ip, nil
}

func (c *Context[request, response]) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *Context[request, response]) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Context[request, response]) Err() error {
	return c.ctx.Err()
}

func (c *Context[request, response]) Value(key any) any {
	return c.ctx.Value(key)
}

func NewContext[request any, response any](ctx context.Context, req request, resp response) *Context[request, response] {
	return &Context[request, response]{
		idx:  -1,
		ctx:  trace.FromContext(ctx),
		hc:   NewChain[request, response](),
		req:  req,
		resp: resp,
		err:  nil,
	}
}
