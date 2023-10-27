package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"git.sr.ht/~primalmotion/simplai/chain"
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

			in := prompt.NewInput(strings.TrimPrefix(input, "/summarize "), nil)
			prmpt, err = summarizer.Format(in)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case strings.HasPrefix(input, "/story "):

			in := prompt.NewInput(strings.TrimPrefix(input, "/story "), nil)
			prmpt, err = storyTeller.Format(in)
			if err != nil {
				fmt.Println(err)
				continue
			}

		default:
			node1 := chain.NewNode(llm, storyTeller)
			node2 := chain.NewNode(llm, summarizer)
			node1.Add(node2)
			fmt.Println(node1.Execute(prompt.NewInput(input, nil)))
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
