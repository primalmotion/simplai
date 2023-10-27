package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"git.sr.ht/~primalmotion/simplai/chain"
	"git.sr.ht/~primalmotion/simplai/node"
	"git.sr.ht/~primalmotion/simplai/prompt"
	"git.sr.ht/~primalmotion/simplai/prompt/storyteller"
	"git.sr.ht/~primalmotion/simplai/prompt/summarizer"
	"git.sr.ht/~primalmotion/simplai/vllm"
)

func main() {

	llm := vllm.NewVLLM(
		"http://cruncher.lan:8000/v1",
		"HuggingFaceH4/zephyr-7b-alpha",
		0.0,
	)

	summarizer := summarizer.NewSummarizer()
	storyTeller := storyteller.NewStoryTeller()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {

		input := strings.TrimSpace(scanner.Text())
		var prmpt string
		var err error

		switch {

		case strings.HasPrefix(input, "/summarize "):

			in := prompt.NewInput(strings.TrimPrefix(input, "/summarize "))
			prmpt, err = summarizer.Format(in)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case strings.HasPrefix(input, "/story "):

			in := prompt.NewInput(strings.TrimPrefix(input, "/story "))
			prmpt, err = storyTeller.Format(in)
			if err != nil {
				fmt.Println(err)
				continue
			}

		default:
			c := chain.New(
				node.New(llm, storyTeller),
				node.New(llm, summarizer),
			)
			fmt.Println(c.Execute(prompt.NewInput(input)))
			continue
		}

		output, err := llm.Infer(prmpt)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(output)
		fmt.Print("> ")
	}
}
