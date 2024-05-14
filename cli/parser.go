package cli

import (
	"github.com/joeycumines/go-prompt"
)

type Parser interface {
	Matcher() prompt.Completer
}

type parserImpl struct {
}
