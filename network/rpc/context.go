package rpc

import (
	"context"
	"github.com/alioth-center/infrastructure/trace"
	"google.golang.org/grpc/peer"
	"net"
	"strings"
	"time"
)

const (
	abortIndex int = 1 << 31
)

type Context[request any, response any] struct {
	idx  int
	ctx  context.Context
	mps  map[string]any
	h    Chain[request, response]
	e    error
	req  request
	resp response
}

func (c *Context[request, response]) Next() {
	c.idx++
	for c.idx < len(c.h) {
		c.h[c.idx](c)
		c.idx++
	}
}

func (c *Context[request, response]) Abort() {
	c.idx = abortIndex
}

func (c *Context[request, response]) IsAborted() bool {
	return c.idx >= abortIndex
}

func (c *Context[request, response]) Set(key string, value any) {
	c.mps[key] = value
}

func (c *Context[request, response]) Get(key string) (value any, exists bool) {
	value, exists = c.mps[key]
	return
}

func (c *Context[request, response]) MustGet(key string, defaultVal any) any {
	if value, exists := c.Get(key); exists {
		return value
	}
	return defaultVal
}

func (c *Context[request, response]) GetString(key string) (value string, exists bool, isString bool) {
	if v, exists := c.Get(key); exists {
		if value, isString = v.(string); isString {
			return value, true, true
		} else {
			return "", true, false
		}
	}

	return "", false, false
}

func (c *Context[request, response]) MustGetString(key string, defaultVal string) string {
	if value, exists, isString := c.GetString(key); exists && isString {
		return value
	}
	return defaultVal
}

func (c *Context[request, response]) GetInt(key string) (value int, exists bool, isInt bool) {
	if v, exists := c.Get(key); exists {
		if value, isInt = v.(int); isInt {
			return value, true, true
		} else {
			return 0, true, false
		}
	}

	return 0, false, false
}

func (c *Context[request, response]) MustGetInt(key string, defaultVal int) int {
	if value, exists, isInt := c.GetInt(key); exists && isInt {
		return value
	}
	return defaultVal
}

func (c *Context[request, response]) GetInt64(key string) (value int64, exists bool, isInt64 bool) {
	if v, exists := c.Get(key); exists {
		if value, isInt64 = v.(int64); isInt64 {
			return value, true, true
		} else {
			return 0, true, false
		}
	}

	return 0, false, false
}

func (c *Context[request, response]) MustGetInt64(key string, defaultVal int64) int64 {
	if value, exists, isInt64 := c.GetInt64(key); exists && isInt64 {
		return value
	}
	return defaultVal
}

func (c *Context[request, response]) GetFloat32(key string) (value float32, exists bool, isFloat32 bool) {
	if v, exists := c.Get(key); exists {
		if value, isFloat32 = v.(float32); isFloat32 {
			return value, true, true
		} else {
			return 0, true, false
		}
	}

	return 0, false, false
}

func (c *Context[request, response]) MustGetFloat32(key string, defaultVal float32) float32 {
	if value, exists, isFloat32 := c.GetFloat32(key); exists && isFloat32 {
		return value
	}
	return defaultVal
}

func (c *Context[request, response]) GetFloat64(key string) (value float64, exists bool, isFloat64 bool) {
	if v, exists := c.Get(key); exists {
		if value, isFloat64 = v.(float64); isFloat64 {
			return value, true, true
		} else {
			return 0, true, false
		}
	}

	return 0, false, false
}

func (c *Context[request, response]) MustGetFloat64(key string, defaultVal float64) float64 {
	if value, exists, isFloat64 := c.GetFloat64(key); exists && isFloat64 {
		return value
	}
	return defaultVal
}

func (c *Context[request, response]) GetBool(key string) (value bool, exists bool, isBool bool) {
	if v, exists := c.Get(key); exists {
		if value, isBool = v.(bool); isBool {
			return value, true, true
		} else {
			return false, true, false
		}
	}

	return false, false, false
}

func (c *Context[request, response]) MustGetBool(key string, defaultVal bool) bool {
	if value, exists, isBool := c.GetBool(key); exists && isBool {
		return value
	}
	return defaultVal
}

func (c *Context[request, response]) GetTime(key string) (value time.Time, exists bool, isTime bool) {
	if v, exists := c.Get(key); exists {
		if value, isTime = v.(time.Time); isTime {
			return value, true, true
		} else {
			return time.Time{}, true, false
		}
	}

	return time.Time{}, false, false
}

func (c *Context[request, response]) MustGetTime(key string, defaultVal time.Time) time.Time {
	if value, exists, isTime := c.GetTime(key); exists && isTime {
		return value
	}
	return defaultVal
}

func (c *Context[request, response]) FunctionGet(key string, fn func(val any)) {
	if v, exists := c.Get(key); exists {
		fn(v)
	}
}

func (c *Context[request, response]) EntryGet(key string, receiver any, marshaller Marshaller) (err error) {
	return marshaller.Unmarshal([]byte(c.MustGetString(key, "")), receiver)
}

func (c *Context[request, response]) EntrySet(key string, payload any, marshaller Marshaller) (err error) {
	var data []byte
	if data, err = marshaller.Marshal(payload); err != nil {
		return err
	}

	c.Set(key, string(data))
	return nil
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
	_, tracedCtx := trace.GetTraceID(ctx)
	c.ctx = tracedCtx
}

func (c *Context[request, response]) TraceID() string {
	traceID, _ := trace.GetTraceID(c.ctx)
	return traceID
}

func (c *Context[request, response]) GetContextClientIP(ctx context.Context) (ip string, err error) {
	if pr, getPrSuccess := peer.FromContext(ctx); !getPrSuccess {
		// 如果获取失败，则返回错误
		return "", NewGetRPCClientIPFailedError()
	} else if pr.Addr == net.Addr(nil) {
		// 如果获取到的是空接口，则返回错误
		return "", NewGetRPCClientIPFailedError()
	} else if pr.Addr.Network() != "tcp" {
		// 如果获取到的不是tcp协议，则返回错误
		return "", NewUnsupportedNetworkError(pr.Addr.Network())
	} else if ipSlice := strings.Split(pr.Addr.String(), ":"); len(ipSlice) != 2 {
		// 如果获取到的不是ip:port(ipv4)格式，则返回错误
		return "", NewInvalidIPAddressError(pr.Addr.String())
	} else {
		// 如果获取到的是ip:port(ipv4)格式，则返回ip
		return ipSlice[0], nil
	}
}

func (c *Context[request, response]) Error() error {
	return c.e
}

func (c *Context[request, response]) SetError(err error) {
	c.e = err
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

func NewContext[request any, response any](ctx context.Context, req request) *Context[request, response] {
	var defaultResp response
	_, tracedCtx := trace.GetTraceID(ctx)
	return &Context[request, response]{
		idx:  -1,
		ctx:  tracedCtx,
		mps:  map[string]any{},
		h:    NewChain[request, response](),
		req:  req,
		resp: defaultResp,
	}
}
