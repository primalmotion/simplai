package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/quick"
	tprompt "github.com/c-bata/go-prompt"
	"github.com/primalmotion/simplai/engine"
	"github.com/primalmotion/simplai/engine/models/mistral"
	"github.com/primalmotion/simplai/engine/ollama"
	"github.com/primalmotion/simplai/engine/openai"
	"github.com/primalmotion/simplai/node"
	"github.com/primalmotion/simplai/prompt"
	"github.com/primalmotion/simplai/tool"
	"github.com/primalmotion/simplai/utils/chunker"
	"github.com/primalmotion/simplai/utils/render"
	"github.com/primalmotion/simplai/vectorstore/memdb"
	"github.com/theckman/yacspin"
)

// huh?
type exitPanic int

func completer(d tprompt.Document) []tprompt.Suggest {

	if w := d.GetWordBeforeCursor(); w != ":" && w != "/" {
		return nil
	}

	s := []tprompt.Suggest{
		{Text: ":flush", Description: "flush the memory of the bot"},
		{Text: ":debug", Description: "turn on and off debug traces"},
		{Text: ":quit", Description: "quit"},
		{Text: "/c", Description: "force conversation tool"},
		{Text: "/C", Description: "force classification tool"},
		{Text: "/r", Description: "force rag tool"},
		{Text: "/R", Description: "force rag with rerank tool"},
		{Text: "/s", Description: "force summarizer tool"},
		{Text: "/S", Description: "force search tool"},
		{Text: "/t", Description: "force story teller tool"},
	}
	return tprompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func matchPrefix(input string, prefix string) (bool, string) {
	if strings.HasPrefix(input, prefix) {
		return true, strings.TrimSpace(
			strings.TrimPrefix(
				input,
				prefix,
			),
		)
	}
	return false, ""
}

func updateSpinner(spinner *yacspin.Spinner, message string) node.Node {
	return node.NewFunc(
		node.Info{Name: "spinner"},
		func(ctx context.Context, in node.Input, err node.Node) (string, error) {
			spinner.Message(message + "...")
			return in.Input(), nil
		})
}

func codeHighlighter() node.Node {
	return node.NewFunc(
		node.Info{Name: "syntax-colorizer"},
		func(ctx context.Context, in node.Input, err node.Node) (string, error) {

			buf := &bytes.Buffer{}

			if lex := lexers.Analyse(in.Input()); lex == nil {
				buf.WriteString(in.Input())
			} else {
				if err := quick.Highlight(buf, in.Input(), lex.Config().Name, "terminal256", "gruvbox"); err != nil {
					return "", err
				}
			}
			return buf.String(), nil
		})
}

func randomOutputSwitcher(freq int, output string) node.Node {
	return node.NewFunc(
		node.Info{Name: "random-scrambler"},
		func(ctx context.Context, in node.Input, n node.Node) (string, error) {
			if rand.Intn(freq) == 0 {
				fmt.Println("SCRAMBLED OUTPUT!!", output)
				return output, nil
			}
			return in.Input(), nil
		})
}

// calm the linter
var _ = randomOutputSwitcher

func run(
	ctx context.Context,
	engineName string,
	model string,
	api string,
	searxURL string,
	debug bool,
) error {

	var llmmodel engine.LLM

	// define our default Inference settings
	defaultInferenceConfig := engine.InferenceConfig{
		Temperature:       0,
		RepetitionPenalty: 1.0,
		TopP:              1,
	}

	switch engineName {
	case "openai":
		var err error
		llmmodel, err = openai.New(api, model, openai.OptionDefaultInferenceConfig(defaultInferenceConfig))
		if err != nil {
			return err
		}
	case "ollama":
		var err error
		llmmodel, err = ollama.New(api, model, ollama.OptionDefaultInferenceConfig(defaultInferenceConfig))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown model type")
	}

	var em engine.Embedder
	em, _ = openai.New(api, "sentence-transformers/all-MiniLM-L6-v2")

	var rr engine.Reranker
	rr, _ = openai.New(api, "cross-encoder/ms-marco-MiniLM-L-6-v2")

	store := memdb.New(em)

	text := `
The cat (Felis catus), commonly referred to as the domestic cat or house cat, is the only domesticated species in the family Felidae. Recent advances in archaeology and genetics have shown that the domestication of the cat occurred in the Near East around 7500 BC. It is commonly kept as a house pet and farm cat, but also ranges freely as a feral cat avoiding human contact. It is valued by humans for companionship and its ability to kill vermin. Because of its retractable claws it is adapted to killing small prey like mice and rats. It has a strong flexible body, quick reflexes, sharp teeth, and its night vision and sense of smell are well developed. It is a social species, but a solitary hunter and a crepuscular predator. Cat communication includes vocalizations like meowing, purring, trilling, hissing, growling, and grunting as well as cat body language. It can hear sounds too faint or too high in frequency for human ears, such as those made by small mammals. It also secretes and perceives pheromones.
Female domestic cats can have kittens from spring to late autumn in temperate zones and throughout the year in equatorial regions, with litter sizes often ranging from two to five kittens. Domestic cats are bred and shown at events as registered pedigreed cats, a hobby known as cat fancy. Animal population control of cats may be achieved by spaying and neutering, but their proliferation and the abandonment of pets has resulted in large numbers of feral cats worldwide, contributing to the extinction of bird, mammal and reptile species.
As of 2017, the domestic cat was the second most popular pet in the United States, with 95.6 million cats owned and around 42 million households owning at least one cat. In the United Kingdom, 26% of adults have a cat, with an estimated population of 10.9 million pet cats as of 2020. As of 2021, there were an estimated 220 million owned and 480 million stray cats in the world.


Etymology and naming
The origin of the English word cat, Old English catt, is thought to be the Late Latin word cattus, which was first used at the beginning of the 6th century.[4] The Late Latin word may be derived from an unidentified African language.[5] The Nubian word kaddîska 'wildcat' and Nobiin kadīs are possible sources or cognates.[6] The Nubian word may be a loan from Arabic قَطّ qaṭṭ ~ قِطّ qiṭṭ.[citation needed]
However, it is "equally likely that the forms might derive from an ancient Germanic word, imported into Latin and thence to Greek and to Syriac and Arabic".[7] The word may be derived from Germanic and Northern European languages, and ultimately be borrowed from Uralic, cf. Northern Sámi gáđfi, 'female stoat', and Hungarian hölgy, 'lady, female stoat'; from Proto-Uralic *käďwä, 'female (of a furred animal)'.[8]
The English puss, extended as pussy and pussycat, is attested from the 16th century and may have been introduced from Dutch poes or from Low German puuskatte, related to Swedish kattepus, or Norwegian pus, pusekatt. Similar forms exist in Lithuanian puižė and Irish puisín or puiscín. The etymology of this word is unknown, but it may have arisen from a sound used to attract a cat.[9][10]
A male cat is called a tom or tomcat[11] (or a gib,[12] if neutered). A female is called a queen[13] (or a molly,[14][user-generated source?] if spayed), especially in a cat-breeding context. A juvenile cat is referred to as a kitten. In Early Modern English, the word kitten was interchangeable with the now-obsolete word catling.[15] A group of cats can be referred to as a clowder or a glaring.[16]`

	text2 := `sentence number one is pretty long. sentence number two is the same.

this is a new paragraph. it has sentence number 3.
`
	sp := chunker.NewSimpleTextSplitter(chunker.OptionsChunkSize(150))
	chunks, _ := sp.Chunk(text)
	for _, c := range chunks {
		fmt.Println("---------")
		fmt.Println(c)
		fmt.Println("---------")
	}

	chunks, _ = sp.Chunk(text2)
	for _, c := range chunks {
		fmt.Println("---------")
		fmt.Println(c)
		fmt.Println("---------")
	}

	os.Exit(0)

	// article, err := readability.FromURL("https://en.wikipedia.org/wiki/Cat", 30*time.Second)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // We split this article into paragraphs and then every paragraph into sentences
	// count := 0
	// for i, p := range strings.Split(article.TextContent, "\n\n") {
	//
	// 	trimmed := trim.Output(p)
	// 	if len(trimmed) > 0 {
	//
	// 		count++
	// 		err := store.AddDocument(
	// 			ctx,
	// 			vectorstore.Document{
	// 				ID:      fmt.Sprintf("doc%d", i),
	// 				Content: trimmed,
	// 			},
	// 		)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 	}
	// }
	//
	// fmt.Println("Documents: ", count)

	debugMode := debug

	spinner, err := yacspin.New(
		yacspin.Config{
			Frequency:       100 * time.Millisecond,
			Suffix:          " ",
			CharSet:         yacspin.CharSets[11],
			SuffixAutoColon: true,
			ColorAll:        true,
			Colors:          []string{"fgYellow"},
			StopColors:      []string{"fgGreen"},
		},
	)
	if err != nil {
		return err
	}

	// this one needs state
	// it's an ugly array for now.
	memstorage := []string{}

	summarizerChain := node.NewSubchainWithName(
		"chain:summarizer",
		prompt.NewSummarizer(),
		mistral.NewLLM(llmmodel),
	)

	classifierChain := node.NewSubchainWithName(
		"chain:classifier",
		mistral.NewChatMemory().WithStorage(&memstorage),
		prompt.NewClassifier(
			prompt.SummarizerInfo,
		),
		mistral.NewLLM(llmmodel),
	)

	storytellerChain := node.NewSubchainWithName(
		"chain:storytelling",
		prompt.NewStoryTeller(),
		mistral.NewLLM(llmmodel),
	)

	searxChain := node.NewSubchainWithName(
		"chain:search",
		tool.NewSearx(searxURL),
		prompt.NewSummarizer(),
		mistral.NewLLM(llmmodel),
	)

	ragChain := node.NewSubchainWithName(
		"chain:rag",
		tool.NewRetriever(store, 2),
		prompt.NewRag(),
		mistral.NewLLM(llmmodel),
	)

	rerankChain := node.NewSubchainWithName(
		"chain:rerank",
		tool.NewRetriever(store, 10),
		tool.NewReranker(rr, 2),
		prompt.NewRag(),
		mistral.NewLLM(llmmodel),
	)

	conversationChain := node.NewSubchainWithName(
		"chain:conversation",
		mistral.NewChatMemory().WithStorage(&memstorage),
		prompt.NewConversation(),
		mistral.NewLLM(llmmodel),
	)

	routerChain := node.NewSubchainWithName(
		"chain:router",
		mistral.NewChatMemory().WithStorage(&memstorage),
		updateSpinner(spinner, "classifying"),
		prompt.NewClassifier(
			prompt.SummarizerInfo,
			tool.SearxInfo,
			prompt.ConversationInfo,
			prompt.StoryTellerInfo,
			prompt.CoderInfo,
		),
		mistral.NewLLM(llmmodel),

		// Comment this out to simpulate an error
		// in json generation.
		// randomOutputSwitcher(2, "not-json"),

		updateSpinner(spinner, "routing"),
		node.NewRouter(
			node.Info{Name: "router"},
			node.RouterSimpleDeciderFunc,
			node.NewSubchainWithName(
				prompt.ConversationInfo.Name,
				updateSpinner(spinner, "thinking"),
				prompt.NewConversation(),
				mistral.NewLLM(llmmodel),
			),
			node.NewSubchainWithName(
				prompt.StoryTellerInfo.Name,
				updateSpinner(spinner, "writing"),
				prompt.NewStoryTeller(),
				mistral.NewLLM(llmmodel),
			),
			node.NewSubchainWithName(
				prompt.SummarizerInfo.Name,
				updateSpinner(spinner, "summarizing"),
				prompt.NewSummarizer(),
				mistral.NewLLM(llmmodel),
			),
			node.NewSubchainWithName(
				tool.SearxInfo.Name,
				updateSpinner(spinner, "searching"),
				tool.NewSearx(searxURL),
				prompt.NewSummarizer(),
				mistral.NewLLM(llmmodel),
			),
			node.NewSubchainWithName(
				prompt.CoderInfo.Name,
				updateSpinner(spinner, "coding"),
				prompt.NewCoder(),
				mistral.NewLLM(llmmodel),
				codeHighlighter(),
			),
		),
	)

	promptExecutor := func(input string) {
		input = strings.TrimSpace(input)

		if input == "" {
			return
		}

		var ch node.Node
		var llmInput node.Input

		if ok, _ := matchPrefix(input, ":debug"); ok {
			debugMode = !debugMode
			render.Box(fmt.Sprintf("debug mode: %t", debugMode), "2")
			return
		}

		if ok, _ := matchPrefix(input, ":quit"); ok {
			// yeah.. there is a story here, obviously.
			// see at the bottom of the function.
			panic(exitPanic(0))
		}

		if ok, in := matchPrefix(input, "/s"); ok {
			llmInput = node.NewInput(in)
			ch = summarizerChain
		}

		if ok, in := matchPrefix(input, "/r"); ok {
			llmInput = node.NewInput(in)
			ch = ragChain
		}

		if ok, in := matchPrefix(input, "/R"); ok {
			llmInput = node.NewInput(in)
			ch = rerankChain
		}

		if ok, in := matchPrefix(input, "/C"); ok {
			llmInput = node.NewInput(in)
			ch = classifierChain
		}

		if ok, in := matchPrefix(input, "/t"); ok {
			llmInput = node.NewInput(in)
			ch = storytellerChain
		}

		if ok, in := matchPrefix(input, "/S"); ok {
			llmInput = node.NewInput(in)
			ch = searxChain
		}

		if ok, in := matchPrefix(input, "/c"); ok {
			llmInput = node.NewInput(in)
			ch = conversationChain
		}

		if ch == nil {
			llmInput = node.NewInput(input)
			ch = routerChain
		}

		if !debugMode {
			_ = spinner.Start()
		}

		output, err := ch.Execute(ctx, llmInput.WithDebug(debugMode))
		if !debugMode {
			_ = spinner.Stop()
		}
		if err != nil {
			render.Box(err.Error(), "1")
			return
		}

		render.Box(output, "12")
	}

	defer func() {
		// See a bit below on why
		// this insanity came to existence.
		v, ok := recover().(exitPanic)
		if ok {
			os.Exit(0)
		}
		panic(v)
	}()

	promptHistory := []string{}
	tprompt.New(
		promptExecutor,
		completer,
		tprompt.OptionPrefix("> "),
		tprompt.OptionHistory(promptHistory),
		tprompt.OptionSuggestionBGColor(9),
		tprompt.OptionDescriptionBGColor(9),
		tprompt.OptionDescriptionTextColor(tprompt.DefaultColor),
		tprompt.OptionSelectedSuggestionTextColor(8),
		tprompt.OptionSelectedSuggestionBGColor(tprompt.DefaultColor),
		tprompt.OptionSelectedDescriptionBGColor(tprompt.DefaultColor),
		tprompt.OptionSelectedDescriptionTextColor(tprompt.DefaultColor),
		// this is an insane way to quit the prompt. but it seems
		// there is no other way. I don't even understand why there
		// is no way to cleanly exit..
		// Without this trick, ctrl-c is dead, as the term is no restored.
		// see: https://github.com/c-bata/go-prompt/issues/59
		tprompt.OptionAddKeyBind(tprompt.KeyBind{
			Key: tprompt.ControlC,
			Fn:  func(*tprompt.Buffer) { panic(exitPanic(0)) },
		}),
	).Run()

	return nil
}
