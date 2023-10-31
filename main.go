package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"git.sr.ht/~primalmotion/simplai/llm/models/mistral"
	"git.sr.ht/~primalmotion/simplai/llm/openai"
	"git.sr.ht/~primalmotion/simplai/node"
	"git.sr.ht/~primalmotion/simplai/prompt"
	"git.sr.ht/~primalmotion/simplai/utils/render"
)

func matchPrefix(input string, prefix string) (bool, string) {

	if strings.HasPrefix(input, fmt.Sprintf("%s", prefix)) {
		return true, strings.TrimSpace(
			strings.TrimPrefix(
				input,
				fmt.Sprintf("%s", prefix),
			),
		)
	}

	return false, ""
}

func main() {

	llmmodel := openai.NewOpenAIAPI(
		"http://cruncher.lan:8000/v1",
		"HuggingFaceH4/zephyr-7b-beta",
		0.0,
	)

	debugMode := true

	// this one needs state
	// it's an ugly array for now.
	memstorage := []string{}

	summarizerChain := node.NewChain(
		node.Desc{Name: "chain:summarizer"},
		mistral.NewChatMemory().WithStorage(&memstorage),
		prompt.NewSummarizer(),
		mistral.NewLLM(llmmodel),
	)

	storytellerChain := node.NewChain(
		node.Desc{Name: "chain:storytelling"},
		prompt.NewStoryTeller(),
		mistral.NewLLM(llmmodel),
	)

	searxChain := node.NewChain(
		node.Desc{Name: "chain:search"},
		mistral.NewChatMemory().WithStorage(&memstorage),
		prompt.NewSearxSearch("https://search.inframonde.me"),
		mistral.NewLLM(llmmodel),
	)

	conversationChain := node.NewChain(
		node.Desc{Name: "chain:conversation"},
		mistral.NewChatMemory().WithStorage(&memstorage),
		prompt.NewConversation(),
		mistral.NewLLM(llmmodel),
	)

	routerChain := node.NewChain(
		node.Desc{Name: "chain:root"},
		mistral.NewChatMemory().WithStorage(&memstorage),
		prompt.NewClassifier(
			prompt.StoryTellerDesc,
			prompt.SummarizerDesc,
			prompt.SearxSearchDesc,
		),
		mistral.NewLLM(llmmodel),
		prompt.NewRouter(
			prompt.NewConversation(),
			prompt.NewStoryTeller(),
			prompt.NewSummarizer(),
			prompt.NewSearxSearch("https://search.inframonde.me"),
		),
		mistral.NewLLM(llmmodel),
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
			ch = routerChain
		}

		if ch == nil {
			llmInput = node.NewInput(input)
			ch = conversationChain
			// llmInput = node.NewInput(input)
			// ch = routerChain
		}

		ctx := context.Background()
		output, err := ch.Execute(ctx, llmInput.WithDebug(debugMode))
		if err != nil {
			render.Box(err.Error(), "1")
			continue
		}

		render.Box(output, "12")
		fmt.Print("> ")
	}
}
