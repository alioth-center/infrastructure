package database

import (
	"context"
	"github.com/alioth-center/infrastructure/logger"
	"time"
)

type Options struct {
	DataSource string
	MaxIdle    int
	MaxOpen    int
	MaxLife    time.Duration
	Timeout    time.Duration
	Stdout     string
	Stderr     string
	DebugLog   bool
}

type Database interface {
	Init(options Options) error
	Migrate(models ...any) error
	Has(table string, query string, args ...any) (exist bool, err error)
	Count(table string, query string, args ...any) (count int64, err error)
	GetOne(receiver any, query string, args ...any) error
	GetAll(receiver any, query string, args ...any) error
	GetPage(receiver any, offset, limit int, query string, args ...any) error
	PickOne(receiver any, query string, args ...any) error
	PickAll(receiver any, length int, query string, args ...any) error
	InsertOne(data any) error
	InsertAll(data any) error
	UpdateOne(data any, query string, args ...any) error
	UpdateAll(data any, query string, args ...any) error
	DeleteOne(query string, args ...any) error
	DeleteAll(query string, args ...any) error
	ExecRaw(sql string, args ...any) error
	QueryRaw(receiver any, sql string, args ...any) error
	HasWithCtx(ctx context.Context, table string, query string, args ...any) (exist bool, err error)
	CountWithCtx(ctx context.Context, table string, query string, args ...any) (count int64, err error)
	GetOneWithCtx(ctx context.Context, receiver any, query string, args ...any) error
	GetAllWithCtx(ctx context.Context, receiver any, query string, args ...any) error
	PickOneWithCtx(ctx context.Context, receiver any, query string, args ...any) error
	PickAllWithCtx(ctx context.Context, receiver any, length int, query string, args ...any) error
	InsertOneWithCtx(ctx context.Context, data any) error
	InsertAllWithCtx(ctx context.Context, data any) error
	UpdateOneWithCtx(ctx context.Context, data any, query string, args ...any) error
	UpdateAllWithCtx(ctx context.Context, data any, query string, args ...any) error
	DeleteOneWithCtx(ctx context.Context, query string, args ...any) error
	DeleteAllWithCtx(ctx context.Context, query string, args ...any) error
	ExecRawWithCtx(ctx context.Context, sql string, args ...any) error
	QueryRawWithCtx(ctx context.Context, receiver any, sql string, args ...any) error
	ExtMethods() ExtMethods
	SetLogger(logger logger.Logger)
}
