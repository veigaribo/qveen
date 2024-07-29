package templates

import (
	"bytes"
	"errors"
	"github.com/veigaribo/template"
)

// HACK: Allows reference to the current executing template.
var CurrentTemplate *template.Template

// Does the same as `template` but has no problem resolving a variable
// template name.
func TemplateTmpl(name string, arg any) (string, error) {
	var buffer bytes.Buffer
	err := CurrentTemplate.ExecuteTemplate(&buffer, name, arg)

	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

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
{{tmpl $tmpl $head}}{{range slice .els 1}}{{- $sep -}}
{{tmpl $tmpl .}}{{end}}{{end -}}
{{- end -}}`
