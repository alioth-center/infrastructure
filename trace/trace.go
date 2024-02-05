package trace

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/peer"
	"net"
)

var (
	traceIDKey = "trace_id"
)

// NewContextWithTraceID 从 context.Background 生成一个新的 context，并附加 trace_id
// Deprecated: use NewContext instead
func NewContextWithTraceID() context.Context {
	return NewContext()
}

// NewContext build a new context with trace_id
func NewContext() context.Context {
	_, ctx := AttachTraceID(context.Background())
	return ctx
}

// NewContextWithTid build a new context with existed trace_id
func NewContextWithTid(traceID string) context.Context {
	return context.WithValue(context.Background(), traceIDKey, traceID)
}

// AttachTraceID 为 context 附加 trace_id
func AttachTraceID(ctx context.Context) (traceID string, result context.Context) {
	traceID = uuid.NewString()
	return traceID, context.WithValue(ctx, traceIDKey, traceID)
}

// GetTraceID 从 context 中获取 trace_id，如果获取失败，则重新生成
func GetTraceID(ctx context.Context) (traceID string, result context.Context) {
	traceIDValue := ctx.Value("trace_id")

	// 获取 trace_id 失败，重新生成
	if traceIDValue == nil {
		return AttachTraceID(ctx)
	}

	traceIDString, convertTraceIDSuccess := traceIDValue.(string)

	// 转换 trace_id 失败，重新生成
	if !convertTraceIDSuccess {
		return AttachTraceID(ctx)
	}

	// 转换 trace_id 成功，返回
	return traceIDString, ctx
}

// ForkTracedContext 从已有 trace_id 的 context 中生成一个新的 context，用于在新的 goroutine 中使用
func ForkTracedContext(ctx context.Context) (forked context.Context) {
	traceID, _ := GetTraceID(ctx)
	return context.WithValue(context.Background(), traceIDKey, traceID)
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
