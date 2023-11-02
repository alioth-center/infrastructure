package rpc

import (
	"github.com/alioth-center/infrastructure/logger"
	"strings"
	"time"
)

type Handler[request any, response any] func(ctx *Context[request, response])

func LogInputAndOutput[request any, response any](ctx *Context[request, response], log logger.Logger, service string, handler string) {
	srv := strings.Builder{}
	srv.WriteString(service)
	srv.WriteString("::")
	srv.WriteString(handler)
	srvString := srv.String()
	ip, getIPErr := ctx.GetContextClientIP(ctx)
	if getIPErr != nil {
		ip = "unknown"
	}

	log.Info(logger.NewFields(ctx).WithService(srvString).WithMessage("rpc started with request").WithData(ctx.req).WithField("client_ip", ip))
	start := time.Now()
	ctx.Next()
	costTime := time.Since(start).String()
	if ctx.Error() != nil {
		log.Error(logger.NewFields(ctx).WithService(srvString).WithMessage("rpc finished with error").WithData(ctx.Error().Error()).WithField("cost", costTime).WithField("client_ip", ip))
	} else {
		log.Info(logger.NewFields(ctx).WithService(srvString).WithMessage("rpc finished success with response").WithData(ctx.resp).WithField("cost", costTime).WithField("client_ip", ip))
	}
}
