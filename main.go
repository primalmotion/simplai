package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"git.sr.ht/~primalmotion/simplai/chain"
	"git.sr.ht/~primalmotion/simplai/llm/openai"
	"git.sr.ht/~primalmotion/simplai/node"
	"git.sr.ht/~primalmotion/simplai/prompt"
	"git.sr.ht/~primalmotion/simplai/utils/render"
)

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
	conversationChain := chain.New(
		prompt.NewConversation("ai", "human").WithPreHook(printPreHook),
		node.NewLLM(llmmodel),
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

		switch {

		case strings.HasPrefix(input, ":debug"):
			debugMode = !debugMode
			render.Box(fmt.Sprintf("debug mode: %t", debugMode), "2")
			fmt.Print("> ")
			continue

		case strings.HasPrefix(input, "/s "):
			llmInput = node.NewInput(strings.TrimPrefix(input, "/s "))
			ch = chain.New(
				prompt.NewSummarizer().WithPreHook(printPreHook),
				node.NewLLM(llmmodel),
			)

		case strings.HasPrefix(input, "/t "):
			llmInput = node.NewInput(strings.TrimPrefix(input, "/t "))
			ch = chain.New(
				prompt.NewStoryTeller().WithPreHook(printPreHook),
				node.NewLLM(llmmodel),
			)

		case strings.HasPrefix(input, "/c "):
			llmInput = node.NewInputWithKeys(
				strings.TrimPrefix(input, "/c "),
				map[string]any{
					"story-teller": "write something, invent a story, tell a tale or a lie.",
					"summarize":    "summarize some text, URL or document.",
					"search":       "fetch information from the internet about people, facts or news.",
				},
			)
			ch = chain.New(
				prompt.NewClassifier().WithPreHook(printPreHook),
				node.NewLLM(llmmodel),
			)

		case strings.HasPrefix(input, "/C "):
			llmInput = node.NewInput(strings.TrimPrefix(input, "/C "))
			ch = chain.New(
				prompt.NewStoryTeller().WithPreHook(printPreHook),
				node.NewLLM(llmmodel),
				prompt.NewSummarizer().WithPreHook(printPreHook),
				node.NewLLM(llmmodel),
			)

		default:
			llmInput = node.NewInput(input)
			ch = conversationChain
		}

		output, err := ch.Execute(llmInput)
		if err != nil {
			render.Box(err.Error(), "1")
			continue
		}

		render.Box(output, "12")
		fmt.Print("> ")
	}
}
