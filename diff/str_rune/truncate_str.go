package str_rune

import (
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

var (
	p, _ = template.New("test").Parse("{{ slice .Content  0 25}}")
)

func normalTruncate(s string, maxWords int) string {
	r := []rune(s)
	if len(r) > maxWords {
		return string(r[:maxWords])
	}
	return s
}

func buildTruncate(s string, maxWords int) string {
	b := strings.Builder{}
	b.WriteString(s)
	if b.Len() > maxWords {
		return b.String()[:maxWords]
	}
	return b.String()
}

func decodeTruncate(s string, maxWords int) string {
	l := utf8.RuneCountInString(s)
	if l <= maxWords {
		return s
	}
	n := 0
	s1 := ""
	for _, d := range strings.Split(s, "") {
		n += utf8.RuneCountInString(d)
		if n > maxWords {
			return s1
		}
		s1 += d
	}
	return s1
}

func TruncateByWords(s string, maxWords int) string {
	processedWords := 0
	wordStarted := false
	for i := 0; i < len(s); {
		r, width := utf8.DecodeRuneInString(s[i:])
		if !isSeparator(r) {
			i += width
			wordStarted = true
			continue
		}
		if !wordStarted {
			i += width
			continue
		}

		wordStarted = false
		processedWords++
		if processedWords == maxWords {
			const ending = "..."
			if (i + len(ending)) >= len(s) {
				// Source string ending is shorter than "..."
				return s
			}
			return s[:i] + ending
		}
		i += width
	}
	// Source string contains less words count than maxWords.
	return s
}

func isSeparator(r rune) bool {
	// ASCII alphanumerics and underscore are not separators
	if r <= 0x7F {
		switch {
		case '0' <= r && r <= '9':
			return false
		case 'a' <= r && r <= 'z':
			return false
		case 'A' <= r && r <= 'Z':
			return false
		case r == '_':
			return false
		}
		return true
	}
	// Letters and digits are not separators
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return false
	}
	// Otherwise, all we can do for now is treat spaces as separators.
	return unicode.IsSpace(r)
}
