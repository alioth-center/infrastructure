package cli

import (
	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/utils/concurrency"
	"github.com/alioth-center/infrastructure/utils/values"
)

type Params struct {
	http.Params
}

func (p Params) Set(key, value string) {
	p.Params[key] = value
}

type NodeType = string

type ExampleType = string

const (
	NodeTypeCommand NodeType = "command"
	NodeTypeOption  NodeType = "option"
)

var (
	handlers  = concurrency.NewHashMap[string, Handler](concurrency.HashMapNodeOptionSmallSize)
	injectors = concurrency.NewHashMap[string, Injector](concurrency.HashMapNodeOptionSmallSize)
)

func RegisterHandler(name string, handler Handler) {
	handlers.Set(name, handler)
}

func GetHandler(name string) (handler Handler, ok bool) {
	return handlers.Get(name)
}

func RegisterInjector(name string, injector Injector) {
	injectors.Set(name, injector)
}

func GetInjector(name string) (injector Injector, ok bool) {
	return injectors.Get(name)
}

type HandlerNotFoundError struct {
	displayLanguage string
	CommandPath     string
	HandlerName     string
}

func (e HandlerNotFoundError) Error() string {
	_, description := i18nPacks[i18nErrHandlerNotFound].GetTranslation(e.displayLanguage)
	return values.NewStringTemplate(description, map[string]string{
		"command": e.CommandPath,
		"handler": e.HandlerName,
	}).Parse()
}

type InjectorNotFoundError struct {
	displayLanguage string
	CommandPath     string
	InjectorName    string
}

func (e InjectorNotFoundError) Error() string {
	_, description := i18nPacks[i18nErrInjectorNotFound].GetTranslation(e.displayLanguage)
	return values.NewStringTemplate(description, map[string]string{
		"command":  e.CommandPath,
		"injector": e.InjectorName,
	}).Parse()
}

const (
	FunctionNameDefaultVersion = "ac-default-version-fn"
	FunctionNameDefaultExit    = "ac-default-exit-fn"
)
