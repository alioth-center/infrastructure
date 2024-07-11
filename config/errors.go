package config

import "errors"

var (
	ErrUnSupportedConfigExtension = errors.New("unsupported config extension")
	ErrConfigFilePathIsDir        = errors.New("config file path is a directory")
	ErrConfigContentNotExists     = errors.New("config content not exists")
)
