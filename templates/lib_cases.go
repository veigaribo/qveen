package templates

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func UpperCase(str string) string {
	return strings.ToUpper(str)
}

func LowerCase(str string) string {
	return strings.ToLower(str)
}

// `str` should be separated by spaces.
func TitleCase(str string) string {
	var builder strings.Builder
	builder.Grow(len(str))

	shouldUp := true

	for _, r := range str {
		if unicode.IsSpace(r) {
			shouldUp = true
		}

		if shouldUp {
			r = unicode.ToTitle(r)
			shouldUp = false
		}

		builder.WriteRune(r)
	}

	return builder.String()
}

// `str` should be separated by spaces.
func PascalCase(str string) string {
	var builder strings.Builder
	builder.Grow(len(str))

	shouldUp := true

	for _, r := range str {
		if unicode.IsSpace(r) {
			shouldUp = true
			continue
		}

		if shouldUp {
			r = unicode.ToUpper(r)
			shouldUp = false
		}

		builder.WriteRune(r)
	}

	return builder.String()
}

// `str` should be separated by spaces.
func CamelCase(str string) string {
	var builder strings.Builder
	builder.Grow(len(str))

	shouldUp := false

	for _, r := range str {
		if unicode.IsSpace(r) {
			shouldUp = true
			continue
		}

		if shouldUp {
			r = unicode.ToUpper(r)
			shouldUp = false
		}

		builder.WriteRune(r)
	}

	return builder.String()
}

// `str` should be separated by spaces.
func SnakeCase(str string) string {
	var builder strings.Builder
	builder.Grow(len(str))

	for _, r := range str {
		if unicode.IsSpace(r) {
			builder.WriteRune('_')
		} else {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

// `str` should be separated by spaces.
func KebabCase(str string) string {
	var builder strings.Builder
	builder.Grow(len(str))

	for _, r := range str {
		if unicode.IsSpace(r) {
			builder.WriteRune('-')
		} else {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

// `str` should be separated by spaces.
func ConstantCase(str string) string {
	var builder strings.Builder
	builder.Grow(len(str))

	for _, r := range str {
		if unicode.IsSpace(r) {
			builder.WriteRune('_')
		} else {
			builder.WriteRune(unicode.ToUpper(r))
		}
	}

	return builder.String()
}

// `str` should be separated by spaces.
func DotCase(str string) string {
	var builder strings.Builder
	builder.Grow(len(str))

	for _, r := range str {
		if unicode.IsSpace(r) {
			builder.WriteRune('.')
		} else {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

// `str` should be separated by spaces.
func SentenceCase(str string) string {
	var builder strings.Builder
	builder.Grow(len(str))

	head, offset := utf8.DecodeRuneInString(str)
	builder.WriteRune(unicode.ToTitle(head))
	builder.WriteString(str[offset:])

	return builder.String()
}
