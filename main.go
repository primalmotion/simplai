package main

import (
	"fmt"
	"git.sr.ht/~primalmotion/fllm/vllm"
)

func main() {
	llm := vllm.NewVLLM(
		"http://cruncher.lan:8000/v1",
		"HuggingFaceH4/zephyr-7b-alpha",
		0.0,
	)

	fmt.Println(llm.Infer("San francisco is a"))
}
