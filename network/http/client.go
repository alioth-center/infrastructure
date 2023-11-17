package http

import (
	"fmt"
	"net/http"
)

type Client interface {
	ExecuteRequest(req Request) Response
}

type client struct {
	c *http.Client
}

func (c *client) ExecuteRequest(req Request) (resp Response) {
	httpReq, buildReqErr := req.build()
	if buildReqErr != nil {
		return NewResponse(nil, fmt.Errorf("execute request error: %w", buildReqErr))
	}
	httpResp, execErr := c.c.Do(httpReq)
	if execErr != nil {
		return NewResponse(nil, fmt.Errorf("execute request error: %w", execErr))
	}

	return NewResponse(httpResp, nil)
}

func NewClient() Client {
	return &client{
		c: &http.Client{},
	}
}
