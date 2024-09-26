package main

import (
	"embed"

	"github.com/alioth-center/infrastructure/cli"
	"github.com/alioth-center/infrastructure/config"
)

var program cli.CommandLine

//go:embed config.yaml
var embeddings embed.FS

func init() {
	cfg := cli.ApplicationConfig{}
	readErr := config.LoadEmbedConfig(&cfg, embeddings, "config.yaml")
	if readErr != nil {
		panic(readErr)
	}

	initCommands()
	program = cli.NewCli(cfg)
}
