package str_rune

import (
	"strings"
)

func normal(s string) string {
	r := []rune(s)
	if len(r) > 100 {
		return string(r[:100])
	}
	return s
}

func build(s string) string {
	b := strings.Builder{}
	b.WriteString(s)
	if b.Len() > 100 {
	}
}