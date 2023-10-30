package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"git.sr.ht/~primalmotion/simplai/chain"
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

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for scanner.Scan() {

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			fmt.Print("> ")
			continue
		}

		var ch *chain.Chain
		var llmInput node.Input

		if ok, _ := matchPrefix(input, ":debug"); ok {
			debugMode = !debugMode
			render.Box(fmt.Sprintf("debug mode: %t", debugMode), "2")
			fmt.Print("> ")
			continue
		}

		if ok, in := matchPrefix(input, "/s"); ok {
			llmInput = node.NewInput(in)
			ch = chain.New(
				prompt.NewSummarizer().WithPreHook(printPreHook),
				node.NewLLM(llmmodel),
			)
		}

		if ok, in := matchPrefix(input, "/t"); ok {
			llmInput = node.NewInput(in)
			ch = chain.New(
				prompt.NewStoryTeller().WithPreHook(printPreHook),
				node.NewLLM(llmmodel),
			)
		}

		if ok, in := matchPrefix(input, "/S"); ok {
			llmInput = node.NewInput(in)
			ch = chain.New(
				memory,
				prompt.NewSearxSearch(memory, "https://search.inframonde.me").WithPreHook(printPreHook),
				node.NewLLM(llmmodel),
			)
		}

		if ok, in := matchPrefix(input, "/c"); ok {
			llmInput = node.NewInput(in).
				WithKeyValue("story-teller", "write something, invent a story, tell a tale or a lie.").
				WithKeyValue("summarize", "summarize some text, URL or document.").
				WithKeyValue("search", "fetch information from the internet about people, facts or news.")
			ch = chain.New(
				prompt.NewClassifier().WithPreHook(printPreHook),
				node.NewLLM(llmmodel),
			)
		}

		if ok, in := matchPrefix(input, "/C"); ok {
			llmInput = node.NewInput(in)
			ch = chain.New(
				prompt.NewStoryTeller().WithPreHook(printPreHook),
				node.NewLLM(llmmodel),
				prompt.NewSummarizer().WithPreHook(printPreHook),
				node.NewLLM(llmmodel),
			)
		}

		if llmInput == nil || ch == nil {
			llmInput = node.NewInput(input)
			ch = chain.New(
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
