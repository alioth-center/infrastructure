package orm

import (
	"context"
	"github.com/alioth-center/infrastructure/database"
	"gorm.io/gorm"
)

const ExtensionName = "orm"

type Extended interface {
	database.Extended
	ExecuteGormFunction(execFunc func(db *gorm.DB) *gorm.DB) error
	QueryGormFunction(receiver any, queryFunc func(db *gorm.DB) *gorm.DB) error
	ExecuteGormTransaction(execFunc func(tx *gorm.DB) error) error
	ExecuteGormFunctionWithCtx(ctx context.Context, execFunc func(db *gorm.DB) *gorm.DB) error
	QueryGormFunctionWithCtx(ctx context.Context, receiver any, queryFunc func(db *gorm.DB) *gorm.DB) error
	ExecuteGormTransactionWithCtx(ctx context.Context, execFunc func(tx *gorm.DB) error) error
}

type extended struct {
	database.Database
	methods database.ExtMethods
}

func (e *extended) ExtensionName() string {
	return ExtensionName
}

func (e *extended) ExecuteGormFunction(execFunc func(db *gorm.DB) *gorm.DB) error {
	return e.methods.Exec(execFunc)
}

func (e *extended) QueryGormFunction(receiver any, queryFunc func(db *gorm.DB) *gorm.DB) error {
	return e.methods.Exec(func(db *gorm.DB) *gorm.DB {
		return queryFunc(db).Scan(receiver)
	})
}

func (e *extended) ExecuteGormTransaction(execFunc func(tx *gorm.DB) error) error {
	return e.methods.Transaction(execFunc)
}

func (e *extended) ExecuteGormFunctionWithCtx(ctx context.Context, execFunc func(db *gorm.DB) *gorm.DB) error {
	return e.methods.ExecCtx(ctx, execFunc)
}

func (e *extended) QueryGormFunctionWithCtx(ctx context.Context, receiver any, queryFunc func(db *gorm.DB) *gorm.DB) error {
	return e.methods.ExecCtx(ctx, func(db *gorm.DB) *gorm.DB {
		return queryFunc(db).Scan(receiver)
	})
}

func (e *extended) ExecuteGormTransactionWithCtx(ctx context.Context, execFunc func(tx *gorm.DB) error) error {
	return e.methods.TransactionCtx(ctx, execFunc)
}
