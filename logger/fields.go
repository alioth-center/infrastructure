package logger

import (
	"context"
	"os"
	"time"

	"github.com/alioth-center/infrastructure/trace"
)

type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
	LevelFatal Level = "fatal"
	LevelPanic Level = "panic"
)

func (l Level) shouldLog(level Level) bool {
	// logger.Level.shouldLog(input.Level)
	return LevelValueMap[l] <= LevelValueMap[level]
}

var (
	workingDirectory = "./"
	timeFormat       = "2006.01.02-15:04:05.000Z07:00"
	LevelValueMap    = map[Level]int{LevelDebug: 0, LevelInfo: 1, LevelWarn: 2, LevelError: 3, LevelFatal: 4, LevelPanic: 5}
)

func init() {
	wd, getWdErr := os.Getwd()
	if getWdErr == nil {
		workingDirectory = wd
	}
}

type Entry struct {
	ctx      context.Context `json:"-" yaml:"-" xml:"-"`
	File     string          `json:"file" yaml:"file" xml:"file"`
	Level    string          `json:"level" yaml:"level" xml:"level"`
	Service  string          `json:"service" yaml:"service" xml:"service"`
	TraceID  string          `json:"trace_id" yaml:"trace_id" xml:"trace_id"`
	CallTime string          `json:"call_time" yaml:"call_time" xml:"call_time"`
	Data     any             `json:"data,omitempty" yaml:"data,omitempty" xml:"data,omitempty"`
	Extra    map[string]any  `json:"extra,omitempty" yaml:"extra,omitempty" xml:"extra,omitempty"`
	Message  string          `json:"message,omitempty" yaml:"message,omitempty" xml:"message,omitempty"`
}

type Fields interface {
	init(ctx context.Context) Fields
	Export() *Entry
	WithTraceID(traceID string) Fields
	WithMessage(message string) Fields
	WithData(data any) Fields
	WithField(key string, value any) Fields
	WithLevel(level Level) Fields
	WithService(service string) Fields
	WithCallTime(callTime time.Time) Fields
	WithBaseFields(base Fields) Fields
}

type fields struct {
	data     any
	level    Level
	file     string
	service  string
	message  string
	callTime string
	extra    map[string]any
	ctx      context.Context
}

// init 初始化日志字段
func (f *fields) init(ctx context.Context) Fields {
	f.file, f.ctx, f.level = trace.Caller(1), trace.FromContext(ctx), LevelInfo

	return f
}

// Export 导出日志字段
func (f *fields) Export() *Entry {
	traceID, ctx := trace.TransformContext(f.ctx)
	f.ctx = ctx

	return &Entry{
		ctx:     f.ctx,
		File:    f.file,
		Level:   string(f.level),
		Service: f.service,
		TraceID: traceID,
		Message: f.message,
		Data:    f.data,
		Extra:   f.extra,
	}
}

// WithTraceID 设置 traceId，如果 traceId 已存在则无效
func (f *fields) WithTraceID(traceId string) Fields {
	if f.ctx.Value(trace.ContextKey()) != nil {
		return f
	}

	f.ctx = context.WithValue(f.ctx, trace.ContextKey(), traceId) // nolint
	return f
}

// WithMessage 设置消息，如果 message 已存在则覆盖
func (f *fields) WithMessage(message string) Fields {
	f.message = message
	return f
}

// WithData 设置数据，如果 data 已存在则覆盖
func (f *fields) WithData(data any) Fields {
	f.data = data
	return f
}

// WithField 设置额外字段，如果 key 已存在则覆盖
func (f *fields) WithField(key string, value any) Fields {
	if f.extra == nil {
		f.extra = map[string]any{}
	}

	f.extra[key] = value
	return f
}

// WithLevel 设置日志级别，不调用则使用 info 级别
func (f *fields) WithLevel(level Level) Fields {
	f.level = level
	return f
}

// WithService 设置服务名称，不调用则使用调用者的函数名
func (f *fields) WithService(service string) Fields {
	f.service = service
	return f
}

// WithCallTime 设置调用时间，不调用则使用当前时间
func (f *fields) WithCallTime(callTime time.Time) Fields {
	if callTime.IsZero() {
		callTime = time.Now()
	}

	f.callTime = callTime.Format(timeFormat)
	return f
}

// WithBaseFields 从基础字段中复制基础字段，需要在调用链的开始处调用，否则会被覆盖
func (f *fields) WithBaseFields(base Fields) Fields {
	entry := base.Export()
	f.level = Level(entry.Level)
	f.service = entry.Service
	f.message = entry.Message
	f.data = entry.Data
	f.extra = entry.Extra
	return f
}

func NewFields(ctx ...context.Context) Fields {
	if len(ctx) == 1 {
		return (&fields{}).init(ctx[0])
	}

	return (&fields{}).init(context.Background())
}

func NewFieldsFromEntry(entry *Entry) Fields {
	return &fields{
		ctx:      trace.NewContextWithTid(entry.TraceID),
		level:    Level(entry.Level),
		file:     entry.File,
		service:  entry.Service,
		message:  entry.Message,
		data:     entry.Data,
		extra:    entry.Extra,
		callTime: entry.CallTime,
	}
}
