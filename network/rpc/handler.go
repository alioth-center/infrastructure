package rpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/alioth-center/infrastructure/cache"
	"github.com/alioth-center/infrastructure/cache/memory"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/trace"
	"github.com/alioth-center/infrastructure/utils/values"
)

var counter cache.Counter = nil

type Handler[request any, response any] func(ctx *Context[request, response])

func getBaseLogFields[request any, response any](ctx *Context[request, response], service, handler string) logger.Fields {
	srv := strings.Builder{}
	srv.WriteString(service)
	srv.WriteString("::")
	srv.WriteString(handler)
	srvString := srv.String()
	ip, getIPErr := ctx.GetContextClientIP()
	if getIPErr != nil {
		ip = "unknown"
	}
	return logger.NewFields(ctx).WithService(srvString).WithField("client_ip", ip)
}

func LogInputAndOutputHandler[request any, response any](log logger.Logger, service string, handler string) func(ctx *Context[request, response]) {
	return func(ctx *Context[request, response]) {
		log.Info(logger.NewFields(ctx).WithBaseFields(getBaseLogFields(ctx, service, handler)).WithMessage("rpc started with request").WithData(ctx.req))
		start := time.Now()
		ctx.Next()
		costTime := time.Since(start).String()
		if ctx.Error() != nil {
			log.Error(logger.NewFields(ctx).WithBaseFields(getBaseLogFields(ctx, service, handler)).WithMessage("rpc finished with error").WithData(ctx.Error().Error()).WithField("cost", costTime))
		} else {
			log.Info(logger.NewFields(ctx).WithBaseFields(getBaseLogFields(ctx, service, handler)).WithMessage("rpc finished success with response").WithData(ctx.resp).WithField("cost", costTime))
		}
	}
}

func RecoveryHandler[request any, response any](log logger.Logger, service string, handler string) func(ctx *Context[request, response]) {
	return func(ctx *Context[request, response]) {
		defer func() {
			if err := recover(); err != nil {
				stack := string(trace.Stack(2))
				fmt.Println(stack)
				log.Panic(logger.NewFields(ctx).WithBaseFields(getBaseLogFields(ctx, service, handler)).WithMessage("rpc panic recovered").WithData(err).WithField("stack", stack))
				ctx.SetResult(values.Nil[response](), err.(error))
			}
		}()
		ctx.Next()
	}
}

func CheckRequestHandler[request any, response any](checkFunc func(request) (bool, string)) func(ctx *Context[request, response]) {
	return func(ctx *Context[request, response]) {
		if success, reason := checkFunc(ctx.req); !success {
			ctx.Abort()
			ctx.SetResult(values.Nil[response](), NewCheckRequestFailedError(reason))
		}
	}
}

func CheckResponseHandler[request any, response any](checkFunc func(response) (bool, string)) func(ctx *Context[request, response]) {
	return func(ctx *Context[request, response]) {
		ctx.Next()
		if success, reason := checkFunc(ctx.resp); !success {
			ctx.Abort()
			ctx.SetResult(values.Nil[response](), NewCheckResponseFailedError(reason))
		}
	}
}

func RequestLimiterHandler[request any, response any](rpd, rpm, rps int) func(ctx *Context[request, response]) {
	if counter == nil {
		counter = memory.NewMemoryCounter(memory.Config{
			EnableInitiativeClean: true,
			CleanIntervalSecond:   180,
			MaxCleanMicroSecond:   100,
			MaxCleanPercentage:    10,
		})
	}

	keyBuilder := func(ip string, limiterType string) string {
		builder := strings.Builder{}
		builder.WriteString(ip)
		builder.WriteString(":")
		builder.WriteString(limiterType)
		return builder.String()
	}

	limiter := func(limitType string, expire time.Duration) func(ip string) bool {
		limit := 0
		switch limitType {
		case "rps":
			limit = rps
		case "rpm":
			limit = rpm
		case "rpd":
			limit = rpd
		}

		if limit <= 0 {
			// 如果限制为0，则不限制
			return func(_ string) bool {
				return true
			}
		}

		return func(ip string) bool {
			key := keyBuilder(ip, limitType)
			rpu := counter.IncreaseWithExpireWhenNotExist(context.Background(), key, 1, expire)
			if rpu == cache.CounterResultEnumFailed || rpu == cache.CounterResultEnumNotEffective {
				// 我们内部计数器出错，默认不限制
				return true
			}
			if rpu.GetValue() > limit {
				// 如果超过限制，则限制
				return false
			}
			return true
		}
	}
	rpsLimiter, rpmLimiter, rpdLimiter := limiter("rps", time.Second), limiter("rpm", time.Minute), limiter("rpd", time.Hour*24)

	return func(ctx *Context[request, response]) {
		ip, getIpErr := ctx.GetContextClientIP()
		if getIpErr != nil {
			ctx.Abort()
			ctx.SetResult(values.Nil[response](), NewGetRPCClientIPFailedError())
			return
		}

		if rpsLimiter(ip) || rpmLimiter(ip) || rpdLimiter(ip) {
			ctx.Abort()
			ctx.SetResult(values.Nil[response](), NewRequestLimiterError("too many requests"))
			return
		}
	}
}

func RequestLimiterHandlerWithFn[request any, response any](rpd, rpm, rps int, _ func(request) (response, error)) func(ctx *Context[request, response]) {
	return RequestLimiterHandler[request, response](rpd, rpm, rps)
}
