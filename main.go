package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"git.sr.ht/~primalmotion/simplai/prompts/storyteller"
	"git.sr.ht/~primalmotion/simplai/prompts/websummarizer"
	"git.sr.ht/~primalmotion/simplai/vllm"
)

func main() {

	llm := vllm.NewVLLM(
		"http://cruncher.lan:8000/v1",
		"HuggingFaceH4/zephyr-7b-alpha",
		0.0,
	)

	summarizer := websummarizer.WebSummarizer{}
	storyTeller := storyteller.StoryTeller{}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {

		input := strings.TrimSpace(scanner.Text())
		var prompt string
		var err error

		switch {

		case strings.HasPrefix(input, "/summarize "):

			prompt, err = summarizer.Format(strings.TrimPrefix(input, "/summarize "))
			if err != nil {
				fmt.Println(err)
				continue
			}

		case strings.HasPrefix(input, "/story "):

			prompt, err = storyTeller.Format(strings.TrimPrefix(input, "/story "))
			if err != nil {
				fmt.Println(err)
				continue
			}

		default:
			prompt = input
		}

		output, err := llm.Infer(prompt)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(output)
		fmt.Print("> ")
	}
}
