package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"git.sr.ht/~primalmotion/simplai/chain"
	"git.sr.ht/~primalmotion/simplai/llm"
	"git.sr.ht/~primalmotion/simplai/node"
	"git.sr.ht/~primalmotion/simplai/prompt"
	"git.sr.ht/~primalmotion/simplai/prompt/classifier"
	"git.sr.ht/~primalmotion/simplai/prompt/storyteller"
	"git.sr.ht/~primalmotion/simplai/prompt/summarizer"
	"git.sr.ht/~primalmotion/simplai/vllm"
)

func main() {

	llmmodel := vllm.NewVLLM(
		"http://cruncher.lan:8000/v1",
		"HuggingFaceH4/zephyr-7b-alpha",
		0.0,
	)

	summarizer := summarizer.NewSummarizer()
	storyTeller := storyteller.NewStoryTeller()
	classifier := classifier.NewClassifier()

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

		case strings.HasPrefix(input, "/classify "):

			in := prompt.NewInputWithKeys(
				strings.TrimPrefix(input, "/classify "),
				map[string]any{
					"story-teller": "The user wants me to invent a story, or a tale or a lie.",
					"summarize":    "The user wants me to summarize some text, or URL or document.",
					"search":       "The user wants me to fetch some information from the internet.",
				},
			)
			prmpt, err = classifier.Format(in)
			if err != nil {
				fmt.Println(err)
				continue
			}

		default:
			c := chain.New(
				node.New(llmmodel, storyTeller),
				node.New(llmmodel, summarizer),
			)
			fmt.Println(c.Execute(prompt.NewInput(input)))
			continue
		}

		// TODO: these options should be part of the prompt or of the node.
		output, err := llmmodel.Infer(prmpt, llm.OptionInferStop("\n", ".", " "))
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("< %s\n", output)
		fmt.Print("> ")
	}
}
