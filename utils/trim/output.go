package trim

import "strings"

// Output Trim extra spaces in a string.
func Output(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
