package caiyun

import (
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/utils/values"
)

const (
	DefaultBaseUrl    = "https://api.caiyunapp.com/"
	DefaultApiVersion = "v2.6/"
	DefaultAdSetUrl   = "https://docs.caiyunapp.com/weather-api/20220518-adcode.csv"
)

type Config struct {
	ApiKey     string
	ApiVersion string
	BaseUrl    string
}

func (cfg Config) AttachEndpoint(location Location, builder http.RequestBuilder) http.RequestBuilder {
	loc := values.BuildStrings("/", values.Float64ToString(location.Longitude()), ",", values.Float64ToString(location.Latitude()))
	baseUrl := values.BuildStrings(cfg.BaseUrl, cfg.ApiVersion, cfg.ApiKey, loc, "/weather")
	return builder.WithPath(baseUrl)
}

func NewConfig(cfg Config) Config {
	if cfg.BaseUrl == "" {
		cfg.BaseUrl = DefaultBaseUrl
	}
	if cfg.ApiVersion == "" {
		cfg.ApiVersion = DefaultApiVersion
	}

	return cfg
}
