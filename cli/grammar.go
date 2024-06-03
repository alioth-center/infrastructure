package cli

import (
	"strings"

	"github.com/alioth-center/infrastructure/network/http"
	"github.com/alioth-center/infrastructure/utils/values"
	"github.com/joeycumines/go-prompt"
)

type TranslatedItem struct {
	Language           string `json:"language" yaml:"language"`
	DisplayKey         string `json:"display_key" yaml:"display_key"`
	DisplayDescription string `json:"display_description" yaml:"display_description"`
}

type TranslationSet struct {
	translations map[string]TranslatedItem
}

func (ts *TranslationSet) GetTranslation(languages ...string) (key, description string) {
	// no translations, cannot return anything
	if ts == nil || len(ts.translations) == 0 {
		return "", ""
	}

	// return the translation if it exists
	for _, language := range languages {
		item, translated := ts.translations[language]
		if translated {
			return item.DisplayKey, item.DisplayDescription
		}

	}

	// return the first translation if no translation exists
	for _, item := range ts.translations {
		key, description = item.DisplayKey, item.DisplayDescription
		break
	}

	return key, description
}

func (ts *TranslationSet) InitTranslations(items []TranslatedItem) {
	ts.translations = map[string]TranslatedItem{}
	for _, item := range items {
		ts.translations[item.Language] = item
	}
}

type Input struct {
	Params   http.Params
	FullText string
	Language []string
}

type grammarContext struct {
	index      int
	fullText   string
	lastWord   string
	splits     []string
	params     http.Params
	ignoreCase bool
	language   []string
}

func fromGrammarContext(ctx *grammarContext) (c *Input) {
	return &Input{
		Params:   ctx.params,
		FullText: ctx.fullText,
		Language: ctx.language,
	}
}

func (gc *grammarContext) next() *grammarContext {
	gc.index++
	return gc
}

func (gc *grammarContext) end() bool {
	return gc.index >= (len(gc.splits) - 1)
}

func newContext(text string, wantLanguage []string, ignoreCase bool) *grammarContext {
	splits := strings.Split(text, " ")
	return &grammarContext{
		index:      0,
		fullText:   text,
		splits:     splits,
		lastWord:   values.LastOfArray(splits),
		ignoreCase: ignoreCase,
		language:   wantLanguage,
		params:     map[string]string{},
	}
}

type grammarNode struct {
	displayKey   string
	nodeType     string
	children     []*grammarNode
	descriptions *TranslationSet
	injector     Injector
	handler      Handler
}

func (gn *grammarNode) matched(key string, strict bool, ignoreCase bool) bool {
	if gn.nodeType == NodeTypeOption {
		return true
	}

	self, word := gn.displayKey, key
	if ignoreCase {
		self, word = strings.ToLower(gn.displayKey), strings.ToLower(key)
	}

	if !strict {
		return strings.HasPrefix(self, word) || key == ""
	}

	return gn.displayKey == key
}

func (gn *grammarNode) indexEndpoints(ctx *grammarContext) (nodes []*grammarNode) {
	if gn.nodeType == NodeTypeOption {
		ctx.params[gn.displayKey] = ctx.splits[ctx.index-1]
	}

	if ctx.end() {
		// end of the line, return matched children
		matched := values.FilterArray(gn.children, func(g *grammarNode) bool {
			return g.matched(ctx.lastWord, false, ctx.ignoreCase)
		})

		if len(matched) > 1 {
			return values.FilterArray(matched, func(g *grammarNode) bool {
				return g.nodeType == NodeTypeCommand
			})
		}

		return matched
	}

	// not the end of the line, continue to the next node
	for _, node := range gn.children {
		if node.matched(ctx.splits[ctx.index], true, ctx.ignoreCase) {
			return node.indexEndpoints(ctx.next())
		}
	}

	return nodes
}

func (gn *grammarNode) indexSuggestions(ctx *grammarContext) (suggestions []prompt.Suggest) {
	eps := gn.indexEndpoints(ctx)
	if len(eps) == 0 {
		return generateErrorPrompt(ctx.language, i18nBadCommand, map[string]string{"command": strings.Join(ctx.splits[:ctx.index+1], " ")})
	}

	for _, node := range eps {
		suggestions = append(suggestions, node.prompts(ctx)...)
	}

	return suggestions
}

func (gn *grammarNode) execute(ctx *grammarContext) {
	children := gn.indexEndpoints(ctx)
	if len(children) == 1 && children[0].matched(ctx.lastWord, true, ctx.ignoreCase) {
		if children[0].nodeType == NodeTypeOption {
			ctx.params[children[0].displayKey] = ctx.lastWord
		}

		if children[0].handler != nil {
			children[0].handler(fromGrammarContext(ctx))
			return
		}

		DefaultHandler()(fromGrammarContext(ctx))
		return
	}

	NoCommandHandler()(fromGrammarContext(ctx))
}

func (gn *grammarNode) prompts(ctx *grammarContext) []prompt.Suggest {
	key, description := gn.descriptions.GetTranslation(ctx.language...)
	if gn.nodeType != NodeTypeOption {
		return []prompt.Suggest{{Text: gn.displayKey, Description: description}}
	}

	if gn.injector != nil {
		var suggestions []prompt.Suggest
		result := gn.injector(fromGrammarContext(ctx))
		for _, item := range result {
			itemKey, itemValue := item.GetTranslation(ctx.language...)
			suggestions = append(suggestions, prompt.Suggest{Text: itemKey, Description: itemValue})
		}

		return suggestions
	}

	return []prompt.Suggest{{Text: key, Description: description}}
}

func (gn *grammarNode) initialize(cfg map[string]CommandConfig, prefixes []string, errDisplayLanguages []string) (errs []error) {
	for name, subCfg := range cfg {
		if subCfg.Type != NodeTypeOption {
			subCfg.Type = NodeTypeCommand
		}

		subNode := &grammarNode{
			displayKey: name,
			nodeType:   subCfg.Type,
			children:   []*grammarNode{},
		}

		if subCfg.Handler != "" {
			handler, got := GetHandler(subCfg.Handler)
			if !got {
				errs = append(errs, HandlerNotFoundError{
					displayLanguages: errDisplayLanguages,
					CommandPath:      strings.Join(append(prefixes, name), " "),
					HandlerName:      subCfg.Handler,
				})
				handler = DefaultHandler()
			}

			subNode.handler = handler
		}

		if subCfg.Type == NodeTypeOption && subCfg.Examples != "" {
			injector, got := GetInjector(subCfg.Examples)
			if !got {
				errs = append(errs, InjectorNotFoundError{
					displayLanguages: errDisplayLanguages,
					CommandPath:      strings.Join(append(prefixes, name), " "),
					InjectorName:     subCfg.Examples,
				})
				injector = nil
			}

			subNode.injector = injector
		}

		if subNode.injector == nil && len(subCfg.Descriptions) > 0 {
			set, items := &TranslationSet{}, make([]TranslatedItem, 0, len(subCfg.Descriptions))
			for _, description := range subCfg.Descriptions {
				items = append(items, TranslatedItem{
					Language:           description.Language,
					DisplayKey:         description.Name,
					DisplayDescription: description.Text,
				})
			}
			set.InitTranslations(items)

			subNode.descriptions = set
		}

		if subNode.injector == nil && subNode.descriptions == nil {
			subNode.descriptions = &TranslationSet{
				translations: map[string]TranslatedItem{
					"en-US": {
						DisplayKey:         name,
						DisplayDescription: subCfg.Type,
					},
				},
			}
		}

		if len(subCfg.Commands) > 0 {
			var next []string
			copy(next, prefixes)
			subErrs := subNode.initialize(subCfg.Commands, append(next, name), errDisplayLanguages)
			errs = append(errs, subErrs...)
		}

		gn.children = append(gn.children, subNode)
	}

	return errs
}
