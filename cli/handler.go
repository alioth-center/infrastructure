package cli

import (
	"fmt"
	"os"
)

type Handler func(input *Input)

type Injector func(input *Input) (result []TranslationSet)

func DefaultHandler() Handler {
	return func(input *Input) {
		_, description := i18nPacks[i18nNoImplement].GetTranslation(input.Language...)
		fmt.Println(description)
	}
}

func NoCommandHandler() Handler {
	return func(input *Input) {
		_, description := i18nPacks[i18nNoCommand].GetTranslation(input.Language...)
		fmt.Println(description + ": " + input.FullText)
	}
}

func StaticInjector(statics []TranslationSet) Injector {
	return func(input *Input) (result []TranslationSet) {
		return statics
	}
}

func DefaultExitHandler() Handler {
	return func(input *Input) {
		os.Exit(0)
	}
}

func init() {
	RegisterHandler(FunctionNameDefaultExit, DefaultExitHandler())
}
