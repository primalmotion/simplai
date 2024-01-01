package reorder

// Distribute reorder and inputay of string so the first indexes
// are in the begining and end and last indexes in the middle.
// This is usefull for RAG prompt as LLM tens to pay more attention
// at the begining and end of the prompt.
func Distribute(input []string) []string {

	result := make([]string, len(input))

	first := 0
	last := len(input) - 1

	turn := false
	for _, item := range input {
		turn = !turn
		if first+last > len(input) {
			break
		}
		if turn {
			result[first] = item
			first++
			continue
		}
		result[last] = item
		last--
	}

	return result

}
