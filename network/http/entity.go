package http

import (
	"context"
	"errors"
	"fmt"

	"github.com/alioth-center/infrastructure/trace"
)

type BaseError struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func NewBaseError(code int, message string) *BaseError {
	return &BaseError{
		ErrorCode:    code,
		ErrorMessage: message,
	}
}

func (e BaseError) Error() string {
	return fmt.Sprintf("(%d) %s", e.ErrorCode, e.ErrorMessage)
}

type BaseResponse[T any] interface {
	BindError(err error)
	BindData(data T)
	BindContext(ctx context.Context)
}

func NewBaseResponse[T any](ctx context.Context, data T, err error) BaseResponse[T] {
	response := &baseResponseImpl[T]{}
	response.BindError(err)
	response.BindData(data)
	response.BindContext(ctx)
	return response
}

type baseResponseImpl[T any] struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	RequestID    string `json:"request_id,omitempty"`
	Data         T      `json:"data,omitempty"`
}

func (r *baseResponseImpl[T]) BindError(err error) {
	var baseErr BaseError
	if errors.As(err, &baseErr) {
		r.ErrorCode, r.ErrorMessage = baseErr.ErrorCode, baseErr.ErrorMessage
		return
	} else if err != nil {
		r.ErrorCode, r.ErrorMessage = ErrorCodeInternalErrorOccurred, err.Error()
	}
}

func (r *baseResponseImpl[T]) BindData(data T) {
	r.Data = data
}

func (r *baseResponseImpl[T]) BindContext(ctx context.Context) {
	r.RequestID = trace.GetTid(ctx)
}
