package prompt

type Formatter interface {
	Format(Input) (string, error)
	StopWords() []string
}
