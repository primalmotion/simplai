package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"git.sr.ht/~primalmotion/simplai/llm/openai"
	"git.sr.ht/~primalmotion/simplai/node"
	"git.sr.ht/~primalmotion/simplai/prompt"
	"git.sr.ht/~primalmotion/simplai/utils/render"
)

func matchPrefix(input string, prefix string) (bool, string) {

	if strings.HasPrefix(input, fmt.Sprintf("%s ", prefix)) {
		return true, strings.TrimSpace(
			strings.TrimPrefix(
				input,
				fmt.Sprintf("%s ", prefix),
			),
		)
	}

	return false, ""
}

func main() {

	llmmodel := openai.NewOpenAIAPI(
		"http://cruncher.lan:8000/v1",
		"HuggingFaceH4/zephyr-7b-alpha",
		0.0,
	)

	debugMode := true
	printPreHook := func(n node.Node, in node.Input) (node.Input, error) {
		if debugMode {
			render.Box(fmt.Sprintf("[%s]\n%s", n.Name(), in.Input()), "4")
		}
		return in, nil
	}

	// this one needs state
	memory := node.NewChatMemory(
		"<|system|>",
		"<|assistant|>",
		"<|user|>",
	)

	summarizerChain := node.NewChain(
		prompt.NewSummarizer().WithPreHook(printPreHook),
		node.NewLLM(llmmodel),
	)

	storytellerChain := node.NewChain(
		prompt.NewStoryTeller().WithPreHook(printPreHook),
		node.NewLLM(llmmodel),
	)

	searxChain := node.NewChain(
		memory,
		prompt.NewSearxSearch("https://search.inframonde.me").WithPreHook(printPreHook),
		node.NewLLM(llmmodel),
	)

	routerChain := node.NewChain(
		memory,
		node.NewChain(
			prompt.NewClassifier(
				prompt.NewStoryTeller(),
				prompt.NewSummarizer(),
				prompt.NewSearxSearch("https://search.inframonde.me"),
			).WithPreHook(printPreHook),
			node.NewLLM(llmmodel),
		),

		node.NewChain(
			prompt.NewRouter(
				prompt.NewStoryTeller(),
				prompt.NewSummarizer(),
				prompt.NewSearxSearch("https://search.inframonde.me"),
			).WithPreHook(printPreHook),
			node.NewLLM(llmmodel),
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
			ch = routerChain
		}

		if llmInput == nil || ch == nil {
			llmInput = node.NewInput(input)
			ch = node.NewChain(
				memory,
				prompt.NewConversation(memory).WithPreHook(printPreHook),
				node.NewLLM(llmmodel),
			)
		}

		ctx := context.Background()
		output, err := ch.Execute(ctx, llmInput)
		if err != nil {
			render.Box(err.Error(), "1")
			continue
		}

		render.Box(output, "12")
		fmt.Print("> ")
	}
}
