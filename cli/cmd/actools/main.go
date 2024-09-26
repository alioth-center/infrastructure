package main

import (
	"github.com/joeycumines/go-prompt"
)

func main() {
	program.Execute(
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
