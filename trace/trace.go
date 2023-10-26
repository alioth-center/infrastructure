package trace

import (
	"context"
	"github.com/google/uuid"
)

// NewContextWithTraceID 从 context.Background 生成一个新的 context，并附加 trace_id
func NewContextWithTraceID() context.Context {
	_, ctx := AttachTraceID(context.Background())
	return ctx
}

// AttachTraceID 为 context 附加 trace_id
func AttachTraceID(ctx context.Context) (traceID string, result context.Context) {
	traceID = uuid.NewString()
	return traceID, context.WithValue(ctx, "trace_id", traceID)
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
