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
	"git.sr.ht/~primalmotion/simplai/utils/render"
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

		if input == "" {
			fmt.Print("> ")
			continue
		}

		var prmpt string
		var err error

		switch {

		case strings.HasPrefix(input, "/s "):

			in := prompt.NewInput(strings.TrimPrefix(input, "/s "))
			prmpt, err = summarizer.Format(in)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case strings.HasPrefix(input, "/t "):

			in := prompt.NewInput(strings.TrimPrefix(input, "/t "))
			prmpt, err = storyTeller.Format(in)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case strings.HasPrefix(input, "/c "):

			in := prompt.NewInputWithKeys(
				strings.TrimPrefix(input, "/c "),
				map[string]any{
					"story-teller": "write something, invent a story, tell a tale or a lie.",
					"summarize":    "summarize some text, URL or document.",
					"search":       "fetch information from the internet about people, facts or news.",
				},
			)
			prmpt, err = classifier.Format(in)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case strings.HasPrefix(input, "/chain "):
			c := chain.New(
				node.New(llmmodel, storyTeller),
				node.New(llmmodel, summarizer),
			)
			fmt.Println(c.Execute(prompt.NewInput(input)))

		default:
			render.Box("unknown action.", "1")
			fmt.Print("> ")
			continue
		}

		// TODO: these options should be part of the prompt or of the node.
		opts := []llm.InferenceOption{
			// llm.OptionInferStop("\n", ".", " "),
		}

		render.Box(prmpt, "3")
		output, err := llmmodel.Infer(prmpt, opts...)
		if err != nil {
			fmt.Println(err)
			continue
		}

		render.Box(output, "12")
		fmt.Print("> ")
	}
}
