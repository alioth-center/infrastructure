package trace

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/peer"
	"net"
)

const (
	defaultTraceIDKey = "trace_id"
)

var (
	traceIDKey = defaultTraceIDKey
)

// SetTraceIDKey it will set the key of trace_id in context globally, only set once.
// if you need to call it, make sure to set it before any other operation.
func SetTraceIDKey(key string) {
	if traceIDKey == defaultTraceIDKey {
		// only set once
		traceIDKey = key
	}
}

// TransformContext transform a context to a traced context.
// if the context is already a traced context, return itself and its trace id.
func TransformContext(ctx context.Context) (traceID string, result context.Context) {
	traced := FromContext(ctx)
	tid := GetTid(traced)
	return tid, traced
}

// GetTid get trace if from a traced context, if not a traced context, return empty string.
// example:
//
//	var tid string
//	tid = trace.GetTid(trace.NewContext()) // tid is a full 36 length uuid string
//	tid = trace.GetTid(context.Background()) // tid is empty string
func GetTid(ctx context.Context) string {
	value := ctx.Value(traceIDKey)
	if value == nil {
		return ""
	}

	tid, success := value.(string)
	if !success {
		return ""
	}

	return tid
}

// FromContext if ctx has trace_id, return itself, else return the context with trace_id
// example:
//
//	tracedCtx := trace.FromContext(ctx)
//
// then you can use tracedCtx to do something
func FromContext(ctx context.Context) (traced context.Context) {
	value := ctx.Value(traceIDKey)
	if value != nil {
		return ctx
	}

	return context.WithValue(ctx, traceIDKey, uuid.NewString())
}

// ForkContext create a new traced context from an existed context.
// only the trace id will be copied, if no trace id, a new trace id will be generated.
func ForkContext(ctx context.Context) (forked context.Context) {
	traced := FromContext(ctx)
	tid := GetTid(traced)
	return NewContextWithTid(tid)
}

// ForkContextWithOpts create a new traced context from an existed context which carries the values in fields
func ForkContextWithOpts(ctx context.Context, fields ...string) (forked context.Context) {
	forked = ForkContext(ctx)
	for _, field := range fields {
		forked = context.WithValue(forked, field, ctx.Value(field))
	}

	return forked
}

// NewContext build a new context with trace id.
func NewContext() context.Context {
	return NewContextWithTid(uuid.NewString())
}

// NewContextWithTid build a new context with existed trace id.
func NewContextWithTid(traceID string) context.Context {
	return context.WithValue(context.Background(), traceIDKey, traceID)
}

// AttachTraceID 为 context 附加 trace_id
func AttachTraceID(ctx context.Context) (traceID string, result context.Context) {
	traceID = uuid.NewString()
	return traceID, context.WithValue(ctx, traceIDKey, traceID)
}

// GetTraceID 从 context 中获取 trace_id，如果获取失败，则重新生成
//
// Deprecated: use TransformContext instead
func GetTraceID(ctx context.Context) (traceID string, result context.Context) {
	return TransformContext(ctx)
}

// GetClientIPFromPeer get client ip from a grpc request, if not a grpc request or no client ip, return empty string
// example:
//
//	clientIP := trace.GetClientIPFromPeer(ctx)
//	if clientIP == "" {
//		// not a grpc request or no client ip
//	}
func GetClientIPFromPeer(ctx context.Context) (ip string) {
	peerCtx, convertSuccess := peer.FromContext(ctx)
	if !convertSuccess {
		// not a grpc context, cannot get client ip, return empty string
		return ""
	}

	if peerCtx.Addr == net.Addr(nil) {
		// no address, return empty string
		return ""
	}

	clientIP, _, err := net.SplitHostPort(peerCtx.Addr.String())
	if err != nil {
		// split host port failed, return empty string
		return ""
	}

	return clientIP
}
