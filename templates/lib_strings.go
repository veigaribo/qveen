package templates

import "strings"

// Similar to the template with the same name.
func TemplateJoinFn(els *[]string, sep string) string {
	var builder strings.Builder

	if len(*els) == 0 {
		return ""
	}

	head := (*els)[0]
	builder.WriteString(head)

	for _, item := range (*els)[1:] {
		builder.WriteString(sep)
		builder.WriteString(item)
	}

	return builder.String()
}
