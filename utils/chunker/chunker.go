package chunker

// chunker the interface to split text into chunks.
type chunker interface {
	Chunk(input string) ([]string, error)
}
