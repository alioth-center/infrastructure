package main

import (
	"embed"
	"text/template"

	"github.com/alioth-center/infrastructure/cli"
	"github.com/alioth-center/infrastructure/config"
)

var (
	//go:embed embeddings/*
	embeddings embed.FS

	//go:embed templates/*
	templates embed.FS

	configuration CliBuilderConfig
	cliOptions    cli.ApplicationConfig
	program       cli.CommandLine

	parsedTemplate  = map[string]*template.Template{}
	i18nCollections = map[string]*cli.TranslationSet{}
)

func init() {
	// read the application configuration
	readAppCfgErr := config.LoadEmbedConfig(&configuration, embeddings, "embeddings/config.yml")
	if readAppCfgErr != nil {
		panic(readAppCfgErr)
	}

	// read the cli configuration
	readCliCfgErr := config.LoadEmbedConfig(&cliOptions, embeddings, configuration.CliConfigFileName)
	if readCliCfgErr != nil {
		panic(readCliCfgErr)
	}

	// load translations
	loadTranslations()

	// read and parsing the templates
	loadTemplates()

	// register handlers
	registerHandlers()

	program = cli.NewCli(cliOptions)
}

func registerHandlers() {
	cli.RegisterHandler("InitProject", InitProject)
}

func loadTemplates() {
	mainTemplate, parseMainErr := template.ParseFS(templates, configuration.TemplateFileNames.MainTemplate)
	if parseMainErr != nil {
		panic(parseMainErr)
	}
	parsedTemplate[configuration.TemplateFileNames.MainTemplate] = mainTemplate

	initTemplate, parseInitErr := template.ParseFS(templates, configuration.TemplateFileNames.InitTemplate)
	if parseInitErr != nil {
		panic(parseInitErr)
	}
	parsedTemplate[configuration.TemplateFileNames.InitTemplate] = initTemplate

	initHandlersTemplate, parseInitHandlersErr := template.ParseFS(templates, configuration.TemplateFileNames.InitHandlersTemplate)
	if parseInitHandlersErr != nil {
		panic(parseInitHandlersErr)
	}
	parsedTemplate[configuration.TemplateFileNames.InitHandlersTemplate] = initHandlersTemplate

	initInjectorsTemplate, parseInitInjectorsErr := template.ParseFS(templates, configuration.TemplateFileNames.InitInjectorsTemplate)
	if parseInitInjectorsErr != nil {
		panic(parseInitInjectorsErr)
	}
	parsedTemplate[configuration.TemplateFileNames.InitInjectorsTemplate] = initInjectorsTemplate

	newHandlerTemplate, parseNewHandlerErr := template.ParseFS(templates, configuration.TemplateFileNames.NewHandlerTemplate)
	if parseNewHandlerErr != nil {
		panic(parseNewHandlerErr)
	}
	parsedTemplate[configuration.TemplateFileNames.NewHandlerTemplate] = newHandlerTemplate

	newInjectorTemplate, parseNewInjectorErr := template.ParseFS(templates, configuration.TemplateFileNames.NewInjectorTemplate)
	if parseNewInjectorErr != nil {
		panic(parseNewInjectorErr)
	}
	parsedTemplate[configuration.TemplateFileNames.NewInjectorTemplate] = newInjectorTemplate
}

func loadTranslations() {
	set := Localization{}
	readTranslationErr := config.LoadEmbedConfig(&set, embeddings, configuration.I18nFileName)
	if readTranslationErr != nil {
		panic(readTranslationErr)
	}

	for key, items := range set {
		translations := &cli.TranslationSet{}
		translations.InitTranslations(items)
		i18nCollections[key] = translations
	}
}
