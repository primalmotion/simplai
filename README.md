# Simplai

SimplAI (pronounced Simpl-Hey) is a library to build LLM based applications in
Go. The goal of SimplAI is to only work with remote inference services like VLLM
or Ollama and not deal with the complexity of doing inference and such locally
in order to provide a sane interface that allows to build powerful LLM
applications.

This is still a work in progress, and the API will change, but things are
shaping up and soon unit testing will be added as the API stabbilizes.

## Chain

The root concept of Simplai is to use stacks of nodes. Each node deals with an
input and returns an output and an eventual error.
This was inspired by Langchain (especially the part about the need to have a
sane a simple interface).

Right now the repositoty contains several prompts, but those are bound to be
moved out of this repository, as it should only contain the basis to build
everything else.

## Get Started

To get started you must have a running instance of VLLM or Ollama, preferably
running Mistral or Zephyr. While the connector for VLLM is using OpenAI style
API, it has never been tested with anything OpenAI (and does not support
authentication yet anyways).

Then run:

    make build
    cd cmd/simplai/
    simplai \
      --engine openai \
      --api http://myserver:8000/v1 \
      --model HuggingFaceH4/zephyr-7b-beta \
      --searx-url https://some.searx.instances

More doc will be added over time. This is still really early in the dev, but the
examples already achieves very good results in our limited testing.

## TODO

- [x] base API
- [x] chain system
- [ ] better memory
- [ ] unit tests
- [ ] args for the main binary
- [ ] retry error back propagation
- [ ] doc strings
- [ ] move some parts out of the repository (prompts)
