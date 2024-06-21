package caiyun

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/alioth-center/infrastructure/logger"
	"github.com/alioth-center/infrastructure/network/http"
)

type Client interface {
	GetCustomWeather(ctx context.Context, location Location, options ...CustomWeatherOptions) (response BaseResponse, err error)
}

type customWeatherOptions = url.Values

type CustomWeatherOptions func(opt *customWeatherOptions)

func WithAlert() CustomWeatherOptions {
	return func(opt *customWeatherOptions) {
		opt.Set("alert", "true")
	}
}

func WithDailysteps(days int) CustomWeatherOptions {
	return func(opt *customWeatherOptions) {
		opt.Set("dailysteps", strconv.Itoa(days))
	}
}

func WithHourlySteps(hours int) CustomWeatherOptions {
	return func(opt *customWeatherOptions) {
		opt.Set("hourlysteps", strconv.Itoa(hours))
	}
}

type client struct {
	executor http.Client
	options  Config
}

func (c *client) GetCustomWeather(ctx context.Context, location Location, options ...CustomWeatherOptions) (response BaseResponse, err error) {
	opt := customWeatherOptions{}
	for _, o := range options {
		o(&opt)
	}

	responseParser, execErr := c.executor.ExecuteRequest(
		c.options.AttachEndpoint(
			location,
			http.NewRequestBuilder().WithContext(ctx).WithQueries(&opt).WithUserAgent(http.AliothClient),
		),
	)
	if execErr != nil {
		return BaseResponse{}, fmt.Errorf("failed to execute request: %w", execErr)
	}

	bindErr := responseParser.BindJson(&response)
	if bindErr != nil {
		return BaseResponse{}, fmt.Errorf("failed to bind response: %w", bindErr)
	}

	return response, nil
}

func NewClient(log logger.Logger, options Config) Client {
	return &client{
		executor: http.NewLoggerClient(log),
		options:  NewConfig(options),
	}
}
