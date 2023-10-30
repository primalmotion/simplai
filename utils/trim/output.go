package trim

import "strings"

func Output(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
