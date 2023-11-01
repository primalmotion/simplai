package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"git.sr.ht/~primalmotion/simplai/llm"
	"git.sr.ht/~primalmotion/simplai/llm/models/mistral"
	"git.sr.ht/~primalmotion/simplai/llm/openai"
	"git.sr.ht/~primalmotion/simplai/node"
	"git.sr.ht/~primalmotion/simplai/prompt"
	"git.sr.ht/~primalmotion/simplai/utils/render"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/quick"
	"github.com/theckman/yacspin"
)

func run(engine string, model string, api string, searxURL string) error {

	var llmmodel llm.LLM

	switch engine {
	case "openai":
		llmmodel = openai.NewOpenAIAPI(api, model, 0.0)
	case "ollama":
		// llmmodel = ollama.New()
		return fmt.Errorf("TODO")
	default:
		return fmt.Errorf("unknown model type")
	}

	debugMode := false
	cfg := yacspin.Config{
		Frequency:       100 * time.Millisecond,
		Suffix:          " ",
		CharSet:         yacspin.CharSets[11],
		SuffixAutoColon: true,
		ColorAll:        true,
		Colors:          []string{"fgYellow"},
		StopColors:      []string{"fgGreen"},
	}

	spinner, err := yacspin.New(cfg)
	if err != nil {
		return err
	}

	// this one needs state
	// it's an ugly array for now.
	memstorage := []string{}

	summarizerChain := node.NewChainWithName(
		"chain:summarizer",
		mistral.NewChatMemory().WithStorage(&memstorage),
		prompt.NewSummarizer(),
		mistral.NewLLM(llmmodel),
	)

	storytellerChain := node.NewChainWithName(
		"chain:storytelling",
		prompt.NewStoryTeller(),
		mistral.NewLLM(llmmodel),
	)

	searxChain := node.NewChainWithName(
		"chain:search",
		mistral.NewChatMemory().WithStorage(&memstorage),
		prompt.NewSearxSearch("https://search.inframonde.me"),
		mistral.NewLLM(llmmodel),
	)

	conversationChain := node.NewChainWithName(
		"chain:conversation",
		mistral.NewChatMemory().WithStorage(&memstorage),
		prompt.NewConversation(),
		mistral.NewLLM(llmmodel),
	)

	routerChain := node.NewChainWithName(
		"chain:root",
		mistral.NewChatMemory().WithStorage(&memstorage),
		updateSpinner(spinner, "classifying"),
		prompt.NewClassifier(
			prompt.StoryTellerInfo,
			prompt.SummarizerInfo,
			prompt.SearxSearchInfo,
			prompt.CoderInfo,
		),
		updateSpinner(spinner, "understanding"),
		mistral.NewLLM(llmmodel),
		updateSpinner(spinner, "routing"),
		prompt.NewRouter(
			node.NewChainWithName(
				"conversation",
				updateSpinner(spinner, "thinking"),
				prompt.NewConversation(),
				mistral.NewLLM(llmmodel),
			),
			node.NewChainWithName(
				"storyteller",
				updateSpinner(spinner, "writing story"),
				prompt.NewStoryTeller(),
				mistral.NewLLM(llmmodel),
			),
			node.NewChainWithName(
				"summarizer",
				updateSpinner(spinner, "summarizing"),
				prompt.NewSummarizer(),
				mistral.NewLLM(llmmodel),
			),
			node.NewChainWithName(
				"search",
				updateSpinner(spinner, "searching the web"),
				prompt.NewSearxSearch("https://search.inframonde.me"),
				mistral.NewLLM(llmmodel),
			),
			node.NewChainWithName(
				"coder",
				updateSpinner(spinner, "coding"),
				prompt.NewCoder(),
				mistral.NewLLM(llmmodel),
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
			spinner.Start()
		}

		ctx := context.Background()
		output, err := ch.Execute(ctx, llmInput.WithDebug(debugMode))
		if !debugMode {
			spinner.Stop()
		}
		if err != nil {
			render.Box(err.Error(), "1")
			fmt.Print("> ")
			continue
		}

		buf := &bytes.Buffer{}

		if lex := lexers.Analyse(output); lex == nil {
			buf.WriteString(output)
		} else {
			quick.Highlight(buf, output, lex.Config().Name, "terminal256", "gruvbox")
		}

		render.Box(buf.String(), "12")
		fmt.Print("> ")
	}

	return nil
}
