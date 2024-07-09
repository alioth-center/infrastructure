package trace

import (
	"context"
	"net"

	"github.com/google/uuid"
	"google.golang.org/grpc/peer"
)

const (
	// defaultTraceIDKey is the default key used to store the trace ID in the context.
	defaultTraceIDKey = "trace_id"
)

var traceIDKey = defaultTraceIDKey

// ContextKey returns the current key used to store the trace ID in the context.
func ContextKey() string {
	return traceIDKey
}

// SetTraceIDKey sets the key used to store the trace ID in the context globally.
// This function should be called only once, before any other operations that use the trace ID.
func SetTraceIDKey(key string) {
	if traceIDKey == defaultTraceIDKey {
		traceIDKey = key
	}
}

// TransformContext transforms a context into a traced context, returning the traced context and its trace ID.
// If the context already contains a trace ID, it returns the context unchanged along with the existing trace ID.
func TransformContext(ctx context.Context) (traceID string, result context.Context) {
	traced := FromContext(ctx)
	tid := GetTid(traced)
	return tid, traced
}

// GetTid retrieves the trace ID from a traced context. If the context does not contain a trace ID, it returns an empty string.
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

// FromContext checks if the context contains a trace ID. If it does, it returns the context unchanged.
// If the context does not contain a trace ID, it returns a new context with a generated trace ID.
func FromContext(ctx context.Context) (traced context.Context) {
	value := ctx.Value(traceIDKey)
	if value != nil {
		return ctx
	}

	return context.WithValue(ctx, traceIDKey, uuid.NewString()) // nolint
}

// ForkContext creates a new traced context from an existing context, copying only the trace ID.
// If the original context does not contain a trace ID, a new trace ID is generated.
func ForkContext(ctx context.Context) (forked context.Context) {
	traced := FromContext(ctx)
	tid := GetTid(traced)
	return NewContextWithTid(tid)
}

// ForkContextWithOpts creates a new traced context from an existing context, copying specified fields along with the trace ID.
func ForkContextWithOpts(ctx context.Context, fields ...string) (forked context.Context) {
	forked = ForkContext(ctx)
	for _, field := range fields {
		forked = context.WithValue(forked, field, ctx.Value(field)) // nolint
	}

	return forked
}

// NewContext creates a new context with a generated trace ID.
func NewContext() context.Context {
	return NewContextWithTid(uuid.NewString())
}

// NewContextWithTid creates a new context with the specified trace ID.
func NewContextWithTid(traceID string) context.Context {
	return context.WithValue(context.Background(), traceIDKey, traceID) // nolint
}

// AttachTraceID attaches a newly generated trace ID to the given context, returning the new trace ID and the updated context.
func AttachTraceID(ctx context.Context) (traceID string, result context.Context) {
	traceID = uuid.NewString()
	return traceID, context.WithValue(ctx, traceIDKey, traceID) // nolint
}

// Context creates a new context with the specified trace ID, overwriting any existing trace ID.
func Context(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID) // nolint
}

// GetClientIPFromPeer extracts the client IP address from a gRPC request context.
// If the context does not originate from a gRPC request or the client IP cannot be determined, it returns an empty string.
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
