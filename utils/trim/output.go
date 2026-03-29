package trim

import "strings"

// Output Trim extra spaces in a string.
func Output(s string) string {
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		out = append(out, strings.Join(strings.Fields(l), " "))
	}
	return strings.Join(out, "\n")
}
