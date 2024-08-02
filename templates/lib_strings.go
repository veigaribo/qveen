package templates

import "strings"

// Similar to the template with the same name.
func TemplateJoinFn(els *[]any, sep string) string {
	var builder strings.Builder

	if len(*els) == 0 {
		return ""
	}

	head := (*els)[0].(string)
	builder.WriteString(head)

	for _, item := range (*els)[1:] {
		builder.WriteString(sep)
		builder.WriteString(item.(string))
	}

	return builder.String()
}
