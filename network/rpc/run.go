package rpc

import (
	"context"
)

func ExecuteChain[request any, response any](handlers Chain[request, response], ctx context.Context, req request, resp response) (result response, err error) {
	rpcCtx := NewContext[request, response](ctx, req, resp)
	handlers.Run(rpcCtx)
	return rpcCtx.GetResponse(), rpcCtx.Error()
}
