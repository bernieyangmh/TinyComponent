package str_rune

import (
	"fmt"
	"strings"
)

func StrFormat(s1, s2, s3 string) string {
	return fmt.Sprintf("%s%s%s", s1, s2, s3)
}

func StrPlus(s1, s2, s3 string) string {
	return s1 + s2 + s3
}

func StrBuilder(s1, s2, s3 string) string {
	b := strings.Builder{}
	b.WriteString(s1)
	b.WriteString(s2)
	b.WriteString(s3)
	return b.String()
}
