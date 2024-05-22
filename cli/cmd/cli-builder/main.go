package main

import (
	"github.com/alioth-center/infrastructure/utils/values"
	"github.com/joeycumines/go-prompt"
)

func main() {
	printWorkingDirectory()
	program.Execute(
		prompt.WithPrefixCallback(func() (prefix string) {
			return values.BuildStringsWithJoinIgnoreEmpty(" ", cliOptions.CliPrefix, configuration.ProjectDirectory, "$ ")
		}),
		prompt.WithShowCompletionAtStart(),
		prompt.WithSuggestionBGColor(prompt.Purple),
		prompt.WithDescriptionBGColor(prompt.DarkGray),
		prompt.WithSuggestionTextColor(prompt.White),
		prompt.WithDescriptionTextColor(prompt.White),
		prompt.WithSelectedDescriptionBGColor(prompt.DarkGray),
		prompt.WithSelectedDescriptionTextColor(prompt.White),
		prompt.WithSelectedSuggestionBGColor(prompt.Black),
		prompt.WithSelectedSuggestionTextColor(prompt.White),
	)
}
