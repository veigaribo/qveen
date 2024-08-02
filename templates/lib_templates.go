package templates

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func TemplateErr(reason string) (string, error) {
	return "", errors.New(reason)
}

func TemplateDump(objs ...any) string {
	var head any

	if len(objs) == 0 {
		goto wrapup
	}

	head = objs[0]
	fmt.Fprint(os.Stderr, dump(head))

	for _, obj := range objs[1:] {
		fmt.Fprint(os.Stderr, " ")
		fmt.Fprint(os.Stderr, dump(obj))
	}

wrapup:
	fmt.Fprint(os.Stderr, "\n")
	return ""
}

func TemplateProbe(obj any) any {
	fmt.Fprintln(os.Stderr, dump(obj))
	return obj
}

// Receives a map with `tmpl`, `els`, `sep` and `pre`.
// Invokes the template `tmpl` for each element in `els`, and separates
// them with `sep`. Will add a `pre` before everything else if there is
// at least one element in `els`.
var TemplateJoinT = `
{{- define "join" -}}

{{- if (not .tmpl) -}}{{err "Missing tmpl in join!"}}{{- end -}}
{{- $tmpl := .tmpl -}}
{{- $pre := or .pre "" -}}
{{- $sep := or .sep "\n" -}}

{{if .els}}{{$pre}}{{$head := index .els 0 -}}
{{template $tmpl $head}}{{range slice .els 1}}{{- $sep -}}
{{template $tmpl .}}{{end}}{{end -}}
{{- end -}}`

func dump(x any) string {
	switch val := x.(type) {
	case nil:
		return "null"
	case int:
		return fmt.Sprint(val)
	case float64:
		return fmt.Sprint(val)
	case bool:
		return fmt.Sprint(val)
	case string:
		return fmt.Sprint(strconv.Quote(val))
	case *[]any:
		var builder strings.Builder
		var head any
		builder.WriteRune('[')

		if len(*val) == 0 {
			goto wrapupSlice
		}

		head = (*val)[0]
		builder.WriteString(dump(head))

		for _, item := range (*val)[1:] {
			builder.WriteString(", ")
			builder.WriteString(dump(item))
		}

	wrapupSlice:
		builder.WriteRune(']')
		return builder.String()
	case *map[string]any:
		{
			var builder strings.Builder
			var isFirst bool
			builder.WriteRune('{')

			if len(*val) == 0 {
				goto wrapupMap
			}

			isFirst = true

			for key, value := range *val {
				if !isFirst {
					builder.WriteString(", ")
				} else {
					isFirst = false
				}

				builder.WriteString(strconv.Quote(key))
				builder.WriteString(": ")
				builder.WriteString(dump(value))
			}

		wrapupMap:
			builder.WriteRune('}')
			return builder.String()
		}
	}

	return fmt.Sprintf("Unrecognized<%[1]T>(%[1]q)", x)
}
