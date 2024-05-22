package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/alioth-center/infrastructure/utils/values"

	"github.com/joeycumines/go-prompt"
	pstrings "github.com/joeycumines/go-prompt/strings"
)

type CommandLine interface {
	Execute(options ...prompt.Option)
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

func GetLanguage(preferredLanguage string, mapping map[string][]string) (languages []string) {
	if len(languageCache) > 0 {
		return languageCache
	}

	var original []string
	if mapping == nil {
		mapping = map[string][]string{}
	}

	try1 := os.Getenv("LANG")
	if try1 != "" {
		// got zh_HK.UTF-8, return zh-HK
		original = append(original, strings.ReplaceAll(strings.Split(try1, ".")[0], "_", "-"))
	}

	try2 := os.Getenv("LC_ALL")
	if try2 != "" {
		// got zh_HK.UTF-8, return zh-HK
		original = append(original, strings.ReplaceAll(strings.Split(try2, ".")[0], "_", "-"))
	}

	if preferredLanguage != "" {
		original = append(original, preferredLanguage)
	}

	original = append(original, "en-US")

	// if mapping is defined, append the mapped languages
	for _, result := range original {
		languages = append(languages, result)
		for key, mapped := range mapping {
			if values.ContainsArray(mapped, result) {
				languages = append(languages, key)
			}
		}
	}

	return values.UniqueArray(languages)
}

type cli struct {
	root *grammarNode
	cfg  ApplicationConfig
}

func (c *cli) getLanguage() (results []string) {
	return GetLanguage(c.cfg.PreferredLanguage, c.cfg.LanguageMapping)
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

func (c *cli) Execute(options ...prompt.Option) {
	opts := append([]prompt.Option{
		prompt.WithCompleter(c.Matcher),
		prompt.WithPrefix(c.cfg.CliPrefix),
	}, options...)
	p := prompt.New(func(s string) {
		c.root.execute(newContext(s, c.getLanguage(), !c.cfg.CaseSensitive))
	}, opts...)
	p.Run()
}
