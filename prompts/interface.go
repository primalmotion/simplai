package prompt

type PromptFormatter interface {
	Format(input string) (string, error)
}
