package cli

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/alioth-center/infrastructure/config"
)

func TestGrammarTree(t *testing.T) {
	lk := sync.Mutex{}
	result := &Input{}

	t.Run("CommonCommand", func(t *testing.T) {
		version := &grammarNode{
			displayKey: "version",
			nodeType:   NodeTypeCommand,
			children:   nil,
			descriptions: &TranslationSet{
				translations: map[string]TranslatedItem{
					"en-US": {
						DisplayKey:         "version",
						DisplayDescription: "Show version information",
						Language:           "en-US",
					},
				},
			},
		}
		info := &grammarNode{
			displayKey: "info",
			nodeType:   NodeTypeOption,
			children:   nil,
			descriptions: &TranslationSet{
				translations: map[string]TranslatedItem{
					"en-US": {
						DisplayKey:         "details",
						DisplayDescription: "details of the status",
						Language:           "en-US",
					},
				},
			},
			handler: func(ctx *Input) {
				result = ctx
			},
		}
		status := &grammarNode{
			displayKey: "status",
			nodeType:   NodeTypeOption,
			children:   []*grammarNode{info},
			descriptions: &TranslationSet{
				translations: map[string]TranslatedItem{
					"en-US": {
						DisplayKey:         "name",
						DisplayDescription: "username of yourself",
						Language:           "en-US",
					},
				},
			},
			injector: StaticInjector([]TranslationSet{
				{
					translations: map[string]TranslatedItem{
						"en-US": {
							DisplayKey:         "sb",
							DisplayDescription: "sb is a sb",
						},
					},
				},
				{
					translations: map[string]TranslatedItem{
						"en-US": {
							DisplayKey:         "me",
							DisplayDescription: "me is a me",
						},
					},
				},
			}),
		}
		me := &grammarNode{
			displayKey: "me",
			nodeType:   NodeTypeCommand,
			children:   []*grammarNode{status},
			descriptions: &TranslationSet{
				translations: map[string]TranslatedItem{
					"en-US": {
						DisplayKey:         "me",
						DisplayDescription: "Show me information",
						Language:           "en-US",
					},
				},
			},
		}
		help := &grammarNode{
			displayKey: "help",
			nodeType:   NodeTypeCommand,
			children:   []*grammarNode{version, me},
			descriptions: &TranslationSet{
				translations: map[string]TranslatedItem{
					"en-US": {
						DisplayKey:         "help",
						DisplayDescription: "Show help information",
						Language:           "en-US",
					},
				},
			},
		}
		root := &grammarNode{
			children: []*grammarNode{help},
		}

		t.Run("Found", func(t *testing.T) {
			suggestions := root.indexSuggestions(newContext("help version", "en-US", true))
			if len(suggestions) != 1 {
				t.Errorf("expected 1 suggestion, got %d", len(suggestions))
			}
			if suggestions[0].Text != "version" {
				t.Errorf("expected 'version', got '%s'", suggestions[0].Text)
			}
			if suggestions[0].Description != "Show version information" {
				t.Errorf("expected 'Show version information', got '%s'", suggestions[0].Description)
			}
		})

		t.Run("NotFound", func(t *testing.T) {
			suggestions := root.indexSuggestions(newContext("fuck version", "en-US", true))
			if len(suggestions) != 1 {
				t.Errorf("expected 1 suggestion, got %d", len(suggestions))
			}
			if suggestions[0].Text != "Bad Command" {
				t.Errorf("expected 'Bad Command', got '%s'", suggestions[0].Text)
			}
		})

		t.Run("TranslationFallback", func(t *testing.T) {
			suggestions := root.indexSuggestions(newContext("help version", "zh-CN", true))
			if len(suggestions) != 1 {
				t.Errorf("expected 1 suggestion, got %d", len(suggestions))
			}
			if suggestions[0].Text != "version" {
				t.Errorf("expected 'version', got '%s'", suggestions[0].Text)
			}
			if suggestions[0].Description != "Show version information" {
				t.Errorf("expected 'Show version information', got '%s'", suggestions[0].Description)
			}
		})

		t.Run("MatchMultiple", func(t *testing.T) {
			suggestions := root.indexSuggestions(newContext("help ", "zh-CN", true))
			if len(suggestions) != 2 {
				t.Errorf("expected 2 suggestion, got %d", len(suggestions))
				t.Log(suggestions)
			}
		})

		t.Run("MatchOptions", func(t *testing.T) {
			suggestions := root.indexSuggestions(newContext("help me sb ", "zh-CN", true))
			if len(suggestions) != 1 {
				t.Errorf("expected 1 suggestion, got %d", len(suggestions))
			}
			if suggestions[0].Text != "details" {
				t.Errorf("expected 'details', got '%s'", suggestions[0].Text)
			}
			if suggestions[0].Description != "details of the status" {
				t.Errorf("expected 'details of the status', got '%s'", suggestions[0].Description)
			}
		})

		t.Run("ExecuteCommand", func(t *testing.T) {
			root.execute(newContext("help me sb want", "zh-CN", true))
			if result.fullText != "help me sb want" {
				t.Errorf("expected 'help me sb want', got '%s'", result.fullText)
			}
			if result.params["info"] != "want" {
				t.Errorf("expected 'want', got '%s'", result.params["info"])
			}
			if result.params["status"] != "sb" {
				t.Errorf("expected 'sb', got '%s'", result.params["status"])
			}
		})

		t.Run("ExecuteCommandNotFound", func(t *testing.T) {
			root.execute(newContext("help you be a sb", "zh-CN", true))
		})

		t.Run("IndexInjector", func(t *testing.T) {
			suggestions := root.indexSuggestions(newContext("help me s", "zh-CN", true))
			if len(suggestions) != 2 {
				t.Errorf("expected 2 suggestion, got %d", len(suggestions))
			}
		})
	})

	t.Run("CommandLine", func(t *testing.T) {
		lk.Lock()
		defer lk.Unlock()
		c := ApplicationConfig{}
		_ = config.LoadConfig(&c, "./test.yml")
		t.Run("ParseConfig", func(t *testing.T) {
			NewCli(c)
		})
		t.Run("Execute", func(t *testing.T) {
			_, err := syscall.Open("/dev/tty", syscall.O_RDONLY, 0)
			if !os.IsNotExist(err) && err != nil {
				t.Skip("cannot open /dev/tty")
			}

			go func() {
				NewCli(c).Execute()
				fmt.Fprintf(os.Stdin, "version\n")
				fmt.Fprintf(os.Stdin, "exit\n")
			}()
			time.Sleep(time.Second)
		})
		t.Run("version", func(t *testing.T) {
			c := &cli{}
			c.version(&Input{})
		})
		t.Run("getLanguage", func(t *testing.T) {
			c := &cli{}
			t.Run("LANG", func(t *testing.T) {
				os.Setenv("LANG", "en-US")
				if c.getLanguage() != "en-US" {
					t.Errorf("expected 'en-US', got '%s'", c.getLanguage())
				}
				os.Unsetenv("LANG")
			})
			t.Run("LC_ALL", func(t *testing.T) {
				os.Setenv("LC_ALL", "zh-CN")
				if c.getLanguage() != "zh-CN" {
					t.Errorf("expected 'zh-CN', got '%s'", c.getLanguage())
				}
				os.Unsetenv("LC_ALL")
			})
		})
	})
}

func TestProgressBar(t *testing.T) {
	t.Run("Progress", func(t *testing.T) {
		task := NewCalculateTask("Testing task", 1000, 0)

		go func() {
			for i := 0; i <= task.totalTasks; i++ {
				if i == task.totalTasks/2 {
					task.RefreshName("Halfway")
				}
				time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
				task.RefreshProgress(i)
			}
		}()

		PrintProgress(task)
	})
}
