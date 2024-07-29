package templates

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

var NotSpecialCase unicode.SpecialCase = unicode.SpecialCase{}
var Case unicode.SpecialCase = NotSpecialCase

func TemplateUpperCase(str string) string {
	return strings.ToUpperSpecial(Case, str)
}

func TemplateLowerCase(str string) string {
	return strings.ToLowerSpecial(Case, str)
}

// `str` should be separated by spaces.
func TemplateTitleCase(str string) string {
	var builder strings.Builder
	builder.Grow(len(str))

	shouldUp := true

	for _, r := range str {
		if unicode.IsSpace(r) {
			shouldUp = true
		} else if shouldUp {
			r = Case.ToTitle(r)
			shouldUp = false
		}

		builder.WriteRune(r)
	}

	return builder.String()
}

// `str` should be separated by spaces.
func TemplatePascalCase(str string) string {
	var builder strings.Builder
	builder.Grow(len(str))

	shouldUp := true

	for _, r := range str {
		if unicode.IsSpace(r) {
			shouldUp = true
			continue
		}

		if shouldUp {
			r = Case.ToTitle(r)
			shouldUp = false
		}

		builder.WriteRune(r)
	}

	return builder.String()
}

// `str` should be separated by spaces.
func TemplateCamelCase(str string) string {
	var builder strings.Builder
	builder.Grow(len(str))

	shouldUp := false

	for _, r := range str {
		if unicode.IsSpace(r) {
			shouldUp = true
			continue
		}

		if shouldUp {
			r = Case.ToTitle(r)
			shouldUp = false
		}

		builder.WriteRune(r)
	}

	return builder.String()
}

// `str` should be separated by spaces.
func TemplateSnakeCase(str string) string {
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
func TemplateKebabCase(str string) string {
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
func TemplateConstantCase(str string) string {
	var builder strings.Builder
	builder.Grow(len(str))

	for _, r := range str {
		if unicode.IsSpace(r) {
			builder.WriteRune('_')
		} else {
			builder.WriteRune(Case.ToUpper(r))
		}
	}

	return builder.String()
}

// `str` should be separated by spaces.
func TemplateDotCase(str string) string {
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
func TemplateSentenceCase(str string) string {
	var builder strings.Builder
	builder.Grow(len(str))

	head, offset := utf8.DecodeRuneInString(str)
	builder.WriteRune(Case.ToTitle(head))
	builder.WriteString(str[offset:])

	return builder.String()
}
