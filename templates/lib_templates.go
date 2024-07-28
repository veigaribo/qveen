package templates

// Receives a map with `items` and `preamble`.
// Invokes a template named `line` per line for each item in `items`.
// If there is at least one item, `preamble` will be output before them.
var TemplateJoin = `
{{- define "join" -}}
{{- $preamble := or .preamble "" -}}
{{- $sep := or .sep "\n" -}}

{{if .items}}{{$preamble}}{{$head := index .items 0 -}}
{{template "line" $head}}{{range slice .items 1}}{{- $sep -}}
{{template "line" .}}{{end}}
{{end -}}
{{- end -}}`
