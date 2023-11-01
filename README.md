# Simplai

> THIS IS A WORK IN PROGRESS AND PREVIEW

SimplAI (pronounced Simpl-Hey) is a library to build LLM based applications in
Go. The goal of SimplAI is to only work with remote inference services like:

- [VLLM](https://vllm.readthedocs.io/en/latest/)
- [Ollama](https://ollama.ai/)

This avoids having to deal with the complexity of doing inference and such
locally in order to provide a sane interface that allows to build powerful LLM
applications.

This is still a work in progress, and the API will change, but things are
shaping up and soon unit testing will be added as the API stabilizes.

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
      --searx-url https://some.searxinstance.org

More doc will be added over time. This is still really early in the dev, but the
examples already achieves very good results in our limited testing.

## Components

The main components in the framework are the `node.Node`, `node.Subchain` and
`node.Input`.

### Node

`Node` is the base object of the framework. A `Node` can be chained to another
`Node`. Together they form a chain.

    [prompt:genstory] -> [llm:mistral] -> [prompt:summarize] -> [llm:zephyr]

A chain can contain nested chains:

    [prompt:search] -> [ [prompt:summarize] -> [llm] ] -> [func:format] -> [llm]

A `Node` can be executed by calling its `Execute()` method. The execution is given a
`context.Context` and an `Input`. It returns a string and an eventual error.

This output will then be fed to the next `Node` in the chain. This process continues
until the execution reaches a node with no `Next()` Node. Then the output is
unstacked and returned by the initial executed Node.

### Subchain

A `Subchain` is a Node that holds a separate chain of `[]Node`.

It can be considered as a single `Node`.

    [node] -> [[node]->[node]->[node]] -> [node]

It can be useful to handle a separate set of `Input`, or `llm.Options` for
instance. `Subchains` can also be used in `Router` nodes, that will
execute a certain `Subchain` based on a condition.

    [classify] -> [llm] -> [router] ?-> [summarize] ->[llm1]
                                    ?-> [search] -> [llm2]
                                    ?-> [generate] -> [llm3]

`Subhain` embeds the `BaseNode` and can be used as any other node.

### Input

TODO

## TODO

- [x] base API
- [x] chain system
- [ ] better memory
- [ ] unit tests
- [x] flags for the main binary
- [ ] retry error back propagation
- [ ] doc strings
- [ ] move some parts out of the repository (prompts)
