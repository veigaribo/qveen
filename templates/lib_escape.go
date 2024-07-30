package templates

import (
	"html"
	"strings"
)

// Escapes characters specified in `chars` by preceding them with a
// backslash.
func EscapeBackslash(chars string, str string) string {
	var builder strings.Builder
	builder.Grow(len(str))

	for _, strRune := range str {
		if strings.ContainsRune(chars, strRune) {
			builder.WriteRune('\\')
		}

		builder.WriteRune(strRune)
	}

	return builder.String()
}

// Escapes characters specified in `chars` by doubling them.
func EscapeDouble(chars string, str string) string {
	var builder strings.Builder
	builder.Grow(len(str))

	for _, strRune := range str {
		if strings.ContainsRune(chars, strRune) {
			builder.WriteRune(strRune)
		}

		builder.WriteRune(strRune)
	}

	return builder.String()
}

// Escapes characters that way.
func EscapeHtml(str string) string {
	return html.EscapeString(str)
}

func Replace(from, to, str string) string {
	return strings.ReplaceAll(str, from, to)
}
