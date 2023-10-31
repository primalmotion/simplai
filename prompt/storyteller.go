package prompt

import (
	"git.sr.ht/~primalmotion/simplai/node"
)

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

type StoryTeller struct {
	*node.Prompt
}

func NewStoryTeller() *StoryTeller {
	return &StoryTeller{
		Prompt: node.NewPrompt(storyTellerTemplate).
			WithName("storyteller").
			WithDescription("write something, invent a story, tell a tale or a lie.").(*node.Prompt),
	}
}

func (n *StoryTeller) WithName(name string) node.Node {
	n.Prompt.WithName(name)
	return n
}

func (n *StoryTeller) WithDescription(desc string) node.Node {
	n.Prompt.WithDescription(desc)
	return n
}

func (n *StoryTeller) WithPreHook(h node.PreHook) node.Node {
	n.Prompt.WithPreHook(h)
	return n
}

func (n *StoryTeller) WithPostHook(h node.PostHook) node.Node {
	n.Prompt.WithPostHook(h)
	return n
}
