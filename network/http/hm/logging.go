// Package hm : http middleware functions, such as logging middleware, authentication middleware, etc.
package hm

import (
	"time"

	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/gin-gonic/gin"
)

type TracingRawRequestLog struct {
	Url        string `json:"url"`
	Method     string `json:"method"`
	Proto      string `json:"proto"`
	StatusCode int    `json:"status_code"`
	TimeCost   string `json:"time_cost"`
	ErrMsg     string `json:"error_message,omitempty"`
}

func TracingRawRequestMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		logging := TracingRawRequestLog{
			Url:    ctx.Request.URL.String(),
			Method: ctx.Request.Method,
			Proto:  ctx.Request.Proto,
		}

		ctx.Next()

		// check if framework set error message
		errMsg, success := ctx.Value(http.ErrorContextKey()).(string)
		if success && errMsg != "" {
			logging.ErrMsg = errMsg
		}
		logging.StatusCode = ctx.Writer.Status()
		logging.TimeCost = time.Since(start).String()

		log.Info(logger.NewFields(ctx).WithMessage("request processed").WithData(logging))
	}
}

func TracingRequestMiddleware[request, response any](log logger.Logger, _ request, _ response) http.Handler[request, response] {
	return func(ctx http.Context[request, response]) {
		// logging request
		log.Info(logger.NewFields(ctx).WithMessage("request received").WithData(ctx.Request()).WithField("headers", ctx.HeaderParams()))

		ctx.Next()

		// logging response
		log.Info(logger.NewFields(ctx).WithMessage("request handled").WithData(ctx.Response()).WithField("headers", ctx.ResponseHeaders()))
	}
}
