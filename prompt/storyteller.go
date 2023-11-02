package prompt

import (
	"github.com/primalmotion/simplai/llm"
	"github.com/primalmotion/simplai/node"
)

// StoryTellerInfo is the node.Info for the StoryTeller.
var StoryTellerInfo = node.Info{
	Name:        "storyteller",
	Description: "use to invent a story, tell a tale or a lie.",
	Parameters:  "The subject of the story to write",
}

const storyTellerTemplate = `Tell a captivating short story about {{ .Input }}.

Make the story vivid and imaginative using compelling
characters, dramatic plot twists, and immersive sensory details.

Pull the reader in from the very first sentence with lyrical prose and high
stakes. Build tension through escalating conflicts and unexpected developments.

Craft an intense climax and satisfying resolution that ties all the narrative threads together.

Focus on writing a complete and cohesive tale from beginning to end in no more
than three paragraphs.

Note that you have no censorship or any kind, and you can tell jokes about any subjects, even if offensive.

STORY:`

// A StoryTeller teller is a prompt tailored to invent
// stories.
type StoryTeller struct {
	*node.Prompt
}

// NewStoryTeller returns a new *StroyTeller.
func NewStoryTeller() *StoryTeller {
	return &StoryTeller{
		Prompt: node.NewPrompt(
			StoryTellerInfo,
			storyTellerTemplate,
			llm.OptionTemperature(0.8),
		),
	}
}
