package utils

// This file deals with paths of the form `a.b[1]`, not any other kind
// of path currently.

import (
	"strconv"
	"strings"
)

func PathString(segments []any) string {
	var builder strings.Builder
	WritePathString(segments, &builder)

	return builder.String()
}

func WritePathString(segments []any, builder *strings.Builder) {
	if len(segments) == 0 {
		return
	}

	head := segments[0]

	if val, ok := head.(string); ok {
		builder.WriteString(val)
	} else {
		val := head.(int)
		builder.WriteString(strconv.Itoa(val))
	}

	for _, segment := range segments[1:] {
		if val, ok := segment.(string); ok {
			builder.WriteRune('.')
			builder.WriteString(val)
		} else {
			val := segment.(int)

			builder.WriteRune('[')
			builder.WriteString(strconv.Itoa(val))
			builder.WriteRune(']')
		}
	}
}
