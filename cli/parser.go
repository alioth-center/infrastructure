package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/joeycumines/go-prompt"
	pstrings "github.com/joeycumines/go-prompt/strings"
)

type CommandLine interface {
	Execute()
	Matcher(document prompt.Document) (suggestions []prompt.Suggest, startChar, endChar pstrings.RuneNumber)
}

func NewCli(cfg ApplicationConfig) CommandLine {
	if cfg.CliPrefix == "" {
		cfg.CliPrefix = "> "
	}

	program := &cli{
		root: &grammarNode{},
		cfg:  cfg,
	}
	RegisterHandler(FunctionNameDefaultVersion, program.version)

	errs := program.root.initialize(cfg.Commands, []string{}, program.getLanguage())
	if cfg.Debug && len(errs) > 0 {
		fmt.Println("Errors/Warnings during initialization:")
		for _, err := range errs {
			fmt.Println(err)
		}
	}

	return program
}

type cli struct {
	root *grammarNode
	cfg  ApplicationConfig
}

func (c *cli) getLanguage() string {
	if c.cfg.PreferredLanguage != "" {
		return c.cfg.PreferredLanguage
	}

	try1 := os.Getenv("LANG")
	if try1 != "" {
		return strings.ReplaceAll(strings.Split(try1, ".")[0], "_", "-")
	}

	try2 := os.Getenv("LC_ALL")
	if try2 != "" {
		return strings.ReplaceAll(strings.Split(try2, ".")[0], "_", "-")
	}

	return "en-US"
}

func (c *cli) newContext(document prompt.Document) *grammarContext {
	return newContext(document.Text, c.getLanguage(), !c.cfg.CaseSensitive)
}

func (c *cli) Matcher(document prompt.Document) (suggestions []prompt.Suggest, startChar, endChar pstrings.RuneNumber) {
	return c.root.indexSuggestions(c.newContext(document)),
		document.CurrentRuneIndex() - pstrings.RuneCountInString(document.GetWordBeforeCursor()),
		document.CurrentRuneIndex()
}

func (c *cli) version(_ *Input) {
	fmt.Println(c.cfg.Copyright)
	fmt.Println(c.cfg.Version, "released at", c.cfg.ReleasedAt)
}

func (c *cli) Execute() {
	for {
		c.root.execute(
			newContext(
				prompt.Input(
					prompt.WithCompleter(c.Matcher),
					prompt.WithPrefix(c.cfg.CliPrefix),
				),
				c.getLanguage(),
				!c.cfg.CaseSensitive,
			),
		)
	}
}
