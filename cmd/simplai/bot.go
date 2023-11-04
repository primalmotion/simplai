package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/quick"
	"github.com/primalmotion/simplai/engine"
	"github.com/primalmotion/simplai/engine/models/mistral"
	"github.com/primalmotion/simplai/engine/ollama"
	"github.com/primalmotion/simplai/engine/openai"
	"github.com/primalmotion/simplai/node"
	"github.com/primalmotion/simplai/prompt"
	"github.com/primalmotion/simplai/tool"
	"github.com/primalmotion/simplai/utils/render"
	"github.com/theckman/yacspin"
)

func matchPrefix(input string, prefix string) (bool, string) {
	if strings.HasPrefix(input, prefix) {
		return true, strings.TrimSpace(
			strings.TrimPrefix(
				input,
				prefix,
			),
		)
	}
	return false, ""
}

func updateSpinner(spinner *yacspin.Spinner, message string) node.Node {
	return node.NewFunc(
		node.Info{Name: "spinner"},
		func(ctx context.Context, in node.Input, err node.Node) (string, error) {
			spinner.Message(message + "...")
			return in.Input(), nil
		})
}

func codeHighlighter() node.Node {
	return node.NewFunc(
		node.Info{Name: "syntax-colorizer"},
		func(ctx context.Context, in node.Input, err node.Node) (string, error) {

			buf := &bytes.Buffer{}

			if lex := lexers.Analyse(in.Input()); lex == nil {
				buf.WriteString(in.Input())
			} else {
				if err := quick.Highlight(buf, in.Input(), lex.Config().Name, "terminal256", "gruvbox"); err != nil {
					return "", err
				}
			}
			return buf.String(), nil
		})
}

func randomOutputSwitcher(freq int, output string) node.Node {
	return node.NewFunc(
		node.Info{Name: "random-scrambler"},
		func(ctx context.Context, in node.Input, n node.Node) (string, error) {
			if rand.Intn(freq) == 0 {
				fmt.Println("SCRAMBLED OUTPUT!!", output)
				return output, nil
			}
			return in.Input(), nil
		})
}

// calm the linter
var _ = randomOutputSwitcher

func run(
	ctx context.Context,
	engineName string,
	model string,
	api string,
	searxURL string,
	debug bool,
) error {

	// cdb := chromadb.New("http://127.0.0.1:8000")
	// col, _ := cdb.CreateCollection(ctx, chromadb.CollectionCreate{
	// 	Name:        "test",
	// 	GetOrCreate: true,
	// })
	// store := chromadb.NewChromaStore(cdb, col.ID)
	//
	// err := store.AddDocument(
	// 	ctx,
	// 	vectorstore.Document{
	// 		ID:        "doc10",
	// 		Content:   "Cats sucks ass",
	// 		Embedding: vectorstore.Embedding{0.1, 0.2, 0.3},
	// 		Metadata:  vectorstore.Metadata{"coucou": "cucul"},
	// 	},
	// 	vectorstore.Document{
	// 		ID:        "doc20",
	// 		Content:   "Dogs are great",
	// 		Embedding: vectorstore.Embedding{0.847734, 0.23784640, 0.389175},
	// 		Metadata:  vectorstore.Metadata{"gougou": "gaga"},
	// 	},
	// )
	// if err != nil {
	// 	panic(err)
	// }
	//
	// doc, err := store.SimilaritySearch(
	// 	ctx,
	// 	vectorstore.Embedding{0.1, 0.2, 0.3},
	// 	10,
	// )
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(doc)

	var llmmodel engine.LLM

	// define our default Inference settings
	defaultInferenceConfig := engine.InferenceConfig{
		Temperature:       0,
		RepetitionPenalty: 1.0,
		TopP:              1,
		TopK:              -1,
	}

	switch engineName {
	case "openai":
		var err error
		llmmodel, err = openai.New(api, model, openai.OptionDefaultInferenceConfig(defaultInferenceConfig))
		if err != nil {
			return err
		}
	case "ollama":
		var err error
		llmmodel, err = ollama.New(api, model, ollama.OptionDefaultInferenceConfig(defaultInferenceConfig))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown model type")
	}

	debugMode := debug

	spinner, err := yacspin.New(
		yacspin.Config{
			Frequency:       100 * time.Millisecond,
			Suffix:          " ",
			CharSet:         yacspin.CharSets[11],
			SuffixAutoColon: true,
			ColorAll:        true,
			Colors:          []string{"fgYellow"},
			StopColors:      []string{"fgGreen"},
		},
	)
	if err != nil {
		return err
	}

	// this one needs state
	// it's an ugly array for now.
	memstorage := []string{}

	summarizerChain := node.NewSubchainWithName(
		"chain:summarizer",
		prompt.NewSummarizer(),
		mistral.NewLLM(llmmodel),
	)

	classifierChain := node.NewSubchainWithName(
		"chain:classifier",
		mistral.NewChatMemory().WithStorage(&memstorage),
		prompt.NewClassifier(
			prompt.SummarizerInfo,
		),
		mistral.NewLLM(llmmodel),
	)

	storytellerChain := node.NewSubchainWithName(
		"chain:storytelling",
		prompt.NewStoryTeller(),
		mistral.NewLLM(llmmodel),
	)

	searxChain := node.NewSubchainWithName(
		"chain:search",
		tool.NewSearx(searxURL),
		prompt.NewSummarizer(),
		mistral.NewLLM(llmmodel),
	)

	conversationChain := node.NewSubchainWithName(
		"chain:conversation",
		mistral.NewChatMemory().WithStorage(&memstorage),
		prompt.NewConversation(),
		mistral.NewLLM(llmmodel),
	)

	routerChain := node.NewSubchainWithName(
		"chain:router",
		mistral.NewChatMemory().WithStorage(&memstorage),
		updateSpinner(spinner, "classifying"),
		prompt.NewClassifier(
			prompt.SummarizerInfo,
			tool.SearxInfo,
			prompt.ConversationInfo,
			prompt.StoryTellerInfo,
			prompt.CoderInfo,
		),
		mistral.NewLLM(llmmodel),

		// Comment this out to simpulate an error
		// in json generation.
		// randomOutputSwitcher(2, "not-json"),

		updateSpinner(spinner, "routing"),
		node.NewRouter(
			node.Info{Name: "router"},
			node.RouterSimpleDeciderFunc,
			node.NewSubchainWithName(
				prompt.ConversationInfo.Name,
				updateSpinner(spinner, "thinking"),
				prompt.NewConversation(),
				mistral.NewLLM(llmmodel),
			),
			node.NewSubchainWithName(
				prompt.StoryTellerInfo.Name,
				updateSpinner(spinner, "writing"),
				prompt.NewStoryTeller(),
				mistral.NewLLM(llmmodel),
			),
			node.NewSubchainWithName(
				prompt.SummarizerInfo.Name,
				updateSpinner(spinner, "summarizing"),
				prompt.NewSummarizer(),
				mistral.NewLLM(llmmodel),
			),
			node.NewSubchainWithName(
				tool.SearxInfo.Name,
				updateSpinner(spinner, "searching"),
				tool.NewSearx(searxURL),
				prompt.NewSummarizer(),
				mistral.NewLLM(llmmodel),
			),
			node.NewSubchainWithName(
				prompt.CoderInfo.Name,
				updateSpinner(spinner, "coding"),
				prompt.NewCoder(),
				mistral.NewLLM(llmmodel),
				codeHighlighter(),
			),
		),
	)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			fmt.Print("> ")
			continue
		}

		var ch node.Node
		var llmInput node.Input

		if ok, _ := matchPrefix(input, ":debug"); ok {
			debugMode = !debugMode
			render.Box(fmt.Sprintf("debug mode: %t", debugMode), "2")
			fmt.Print("> ")
			continue
		}

		if ok, in := matchPrefix(input, "/s"); ok {
			llmInput = node.NewInput(in)
			ch = summarizerChain
		}

		if ok, in := matchPrefix(input, "/C"); ok {
			llmInput = node.NewInput(in)
			ch = classifierChain
		}

		if ok, in := matchPrefix(input, "/t"); ok {
			llmInput = node.NewInput(in)
			ch = storytellerChain
		}

		if ok, in := matchPrefix(input, "/S"); ok {
			llmInput = node.NewInput(in)
			ch = searxChain
		}

		if ok, in := matchPrefix(input, "/c"); ok {
			llmInput = node.NewInput(in)
			ch = conversationChain
		}

		if ch == nil {
			llmInput = node.NewInput(input)
			ch = routerChain
		}

		if !debugMode {
			_ = spinner.Start()
		}

		output, err := ch.Execute(ctx, llmInput.WithDebug(debugMode))
		if !debugMode {
			_ = spinner.Stop()
		}
		if err != nil {
			render.Box(err.Error(), "1")
			fmt.Print("> ")
			continue
		}

		render.Box(output, "12")
		fmt.Print("> ")
	}

	return nil
}
