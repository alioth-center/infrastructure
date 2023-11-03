package rpc

import (
	"fmt"
	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/trace"
	"github.com/alioth-center/infrastructure/utils"
	"strings"
	"time"
)

type Handler[request any, response any] func(ctx *Context[request, response])

func getBaseLogFields[request any, response any](ctx *Context[request, response], service, handler string) logger.Fields {
	srv := strings.Builder{}
	srv.WriteString(service)
	srv.WriteString("::")
	srv.WriteString(handler)
	srvString := srv.String()
	ip, getIPErr := ctx.GetContextClientIP(ctx)
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
				ctx.SetResult(utils.NilValue[response](), err.(error))
			}
		}()
		ctx.Next()
	}
}

func CheckRequestHandler[request any, response any](checkFunc func(request) (bool, string)) func(ctx *Context[request, response]) {
	return func(ctx *Context[request, response]) {
		if success, reason := checkFunc(ctx.req); !success {
			ctx.Abort()
			ctx.SetResult(utils.NilValue[response](), NewCheckRequestFailedError(reason))
		}
	}
}

func CheckResponseHandler[request any, response any](checkFunc func(response) (bool, string)) func(ctx *Context[request, response]) {
	return func(ctx *Context[request, response]) {
		ctx.Next()
		if success, reason := checkFunc(ctx.resp); !success {
			ctx.Abort()
			ctx.SetResult(utils.NilValue[response](), NewCheckResponseFailedError(reason))
		}
	}
}
