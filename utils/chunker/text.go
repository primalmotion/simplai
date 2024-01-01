package chunker

import (
	"fmt"
	"strings"

	"github.com/primalmotion/simplai/utils/trim"
)

// SimpleTextChunker is a SimpleTextChunker
type SimpleTextChunker struct {
	Options
}

func NewSimpleTextSplitter(opts ...Option) *SimpleTextChunker {
	options := DefaultOptions()
	for _, o := range opts {
		o(&options)
	}
	return &SimpleTextChunker{Options: options}
}

// Split implement the textsplitter interface.
func (s SimpleTextChunker) Chunk(input string) ([]string, error) {

	overlap := s.ChunkSize * s.ChunkOverlapPercent / 100
	target := s.ChunkSize - overlap

	cleanInput := trim.Output(input)

	// Split our input into Chunks accounting for overlap
	// var output []string
	// we want to crate chunk of size target with overlap
	c := s.split(cleanInput, target, overlap)
	c.Merge(s.ChunkSize)

	// // Expand our Chunks to create overlap
	// output := []string{}
	// for i := range parts {
	// 	var chunk string
	// 	// add a bit from above and below if we can
	// 	if i > 0 {
	// 		for j := 0; j<len(parts) -i; j++ {
	//
	//
	// 		}
	//
	// 	}
	// 	//
	// 	if i < len(parts)-1 {
	//
	// 	}
	//
	// }
	// at this point we have like paragraph / sentence / long string
	// that are chunk size minus the overlap
	// we can then add a bit of above / below to each

	return c.data, nil

}

// chunk represent a chunk
type chunk struct {
	separator string
	data      []string
	chunks    []chunk
}

// Merge will merge back chunks data using the separator for
// a specified count that can be negative to merge the last
// N items or positibe to merge the first N items.
func (c *chunk) Merge(count int) string {
	for _, child := range c.chunks {
		c.data = append(c.data, child.Merge(count))
	}
	var from, to int
	if count < 0 {
		to = len(c.data) - 1
		if -count > len(c.data) {
			from = 0
		} else {
			from = len(c.data) - 1 + count
		}
	} else {
		from = 0
		if count > len(c.data) {
			to = len(c.data) - 1
		} else {
			to = count
		}
	}
	if len(c.data) > 1 {
		return strings.Join(c.data[from:to], c.separator)
	}
	if len(c.data) == 1 {
		return strings.TrimRight(c.data[0]+c.separator, " ")
	}
	return c.separator
}

// split will split a text given a separator. This can recurse
// accros the seperator if the chunk size is not met.
func (s SimpleTextChunker) split(input string, target int, overlap int) chunk {

	// Find our splitter according to the input
	var separator string
	for _, sep := range s.Separators {
		if strings.Contains(input, sep) {
			separator = sep
			break
		}
	}

	c := chunk{
		separator: separator,
	}

	parts := strings.Split(input, separator)

	for idx := range parts {
		var prev, current, next string
		if idx > 0 && overlap != 0 {
			prev = parts[idx-1]
		}
		if idx < len(parts)-1 && overlap != 0 {
			next = parts[idx+1]
		}
		coverlap := overlap
		if prev != "" && next != "" {
			coverlap = overlap / 2
		}
		current = parts[idx]

		if s.countToken(current) < target {
			// add a bit of above
			if prev != "" {
				if s.countToken(prev) > coverlap {
					as := s.split(prev, coverlap, 0)
					fmt.Printf("\n------prev %s\n%s", prev, as.Merge(-overlap))
					c.data = append(c.data, as.Merge(-coverlap))
				} else {
					c.data = append(c.data, prev)
				}
			}
			c.data = append(c.data, current)
			if next != "" {
				if s.countToken(next) > coverlap {
					as := s.split(next, coverlap, 0)
					c.data = append(c.data, as.Merge(coverlap))
				} else {
					c.data = append(c.data, next)
				}
			}
			continue
		}

		//otherwise we split
		c.chunks = append(c.chunks, s.split(current, target, overlap))

	}

	return c
}

// countToken is a dumb word counter. Idealy should be a token counter.
// It will
func (s SimpleTextChunker) countToken(t string) int {
	return int((1.0 + s.TokenRatio) * float32(len(strings.Fields(t))))
}
