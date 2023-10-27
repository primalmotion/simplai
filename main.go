package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"git.sr.ht/~primalmotion/simplai/chain"
	"git.sr.ht/~primalmotion/simplai/node"
	"git.sr.ht/~primalmotion/simplai/prompt"
	"git.sr.ht/~primalmotion/simplai/prompt/classifier"
	"git.sr.ht/~primalmotion/simplai/prompt/storyteller"
	"git.sr.ht/~primalmotion/simplai/prompt/summarizer"
	"git.sr.ht/~primalmotion/simplai/utils/render"
	"git.sr.ht/~primalmotion/simplai/vllm"
)

func main() {

	llmmodel := vllm.NewVLLM(
		"http://cruncher.lan:8000/v1",
		"HuggingFaceH4/zephyr-7b-alpha",
		0.0,
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
		var llmInput prompt.Input

		switch {

		case strings.HasPrefix(input, "/s "):
			llmInput = prompt.NewInput(strings.TrimPrefix(input, "/s "))
			ch = chain.New(
				summarizer.NewSummarizer(),
				node.NewPrintNode(),
				node.NewLLM(llmmodel),
			)

		case strings.HasPrefix(input, "/t "):
			llmInput = prompt.NewInput(strings.TrimPrefix(input, "/t "))
			ch = chain.New(
				storyteller.NewStoryTeller(),
				node.NewPrintNode(),
				node.NewLLM(llmmodel),
			)

		case strings.HasPrefix(input, "/c "):
			llmInput = prompt.NewInputWithKeys(
				strings.TrimPrefix(input, "/c "),
				map[string]any{
					"story-teller": "write something, invent a story, tell a tale or a lie.",
					"summarize":    "summarize some text, URL or document.",
					"search":       "fetch information from the internet about people, facts or news.",
				},
			)
			ch = chain.New(
				classifier.NewClassifier(),
				node.NewPrintNode(),
				node.NewLLM(llmmodel),
			)

		default:
			render.Box("Unknown action.", "1")
			fmt.Print("> ")
			continue
		}

		output, err := ch.Execute(llmInput)
		if err != nil {
			fmt.Println(err)
			continue
		}

		render.Box(output, "12")
		fmt.Print("> ")
	}
}
