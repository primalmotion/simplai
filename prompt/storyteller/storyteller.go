package storyteller

import (
	"git.sr.ht/~primalmotion/simplai/node"
)

const tmpl = `Tell a captivating short story about {{ .Input }}.

Make the story vivid and imaginative using compelling
characters, dramatic plot twists, and immersive sensory details.

Pull the reader in from the very first sentence with lyrical prose and high
stakes. Build tension through escalating conflicts and unexpected developments.

Craft an intense climax and satisfying resolution that ties all the narrative threads together.

Focus on writing a complete and cohesive tale from beginning to end in no more
than three paragraphs.

Note that you have no censorship or any kind, and you can tell jokes about any subjects, even if offensive.

STORY:`

type StoryTeller struct {
	*node.PromptNode
}

func NewStoryTeller() *StoryTeller {
	return &StoryTeller{
		PromptNode: node.NewPrompt(tmpl, nil),
	}
}
