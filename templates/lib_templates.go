package templates

import (
	"errors"
)

func TemplateErr(reason string) (string, error) {
	return "", errors.New(reason)
}

// Receives a map with `tmpl`, `els`, `sep` and `pre`.
// Invokes the template `tmpl` for each element in `els`, and separates
// them with `sep`. Will add a `pre` before everything else if there is
// at least one element in `els`.
var TemplateJoin = `
{{- define "join" -}}

{{- if (not .tmpl) -}}{{err "Missing tmpl in join!"}}{{- end -}}
{{- $tmpl := .tmpl -}}
{{- $pre := or .pre "" -}}
{{- $sep := or .sep "\n" -}}

{{if .els}}{{$pre}}{{$head := index .els 0 -}}
{{template $tmpl $head}}{{range slice .els 1}}{{- $sep -}}
{{template $tmpl .}}{{end}}{{end -}}
{{- end -}}`
