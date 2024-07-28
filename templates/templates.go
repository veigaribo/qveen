package templates

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
)

var Funcs = template.FuncMap{
	"uppercase":    TemplateUpperCase,
	"lowercase":    TemplateLowerCase,
	"titlecase":    TemplateTitleCase,
	"pascalcase":   TemplatePascalCase,
	"camelcase":    TemplateCamelCase,
	"snakecase":    TemplateSnakeCase,
	"kebabcase":    TemplateKebabCase,
	"constcase":    TemplateConstantCase,
	"dotcase":      TemplateDotCase,
	"sentencecase": TemplateSentenceCase,

	"map": TemplateMap,

	"err":  TemplateErr,
	"tmpl": TemplateTmpl,

	"jq1": TemplateJq1,
	"jqn": TemplateJqN,
}

var LeftDelim string = ""
var RightDelim string = ""

var baseTemplate *template.Template

func init() {
	var err error

	baseTemplate = template.
		New("qveen").
		Delims(LeftDelim, RightDelim).
		Funcs(Funcs)

	builtinTemplates := []string{
		TemplateJoin,
	}

	for i, builtin := range builtinTemplates {
		baseTemplate, err = baseTemplate.Parse(builtin)

		if err != nil {
			panic(fmt.Errorf("Failed to parse builtin template #%d! %w", i, err))
		}
	}
}

type WrappedTemplate struct {
	Template *template.Template
}

func (w WrappedTemplate) Parse(text string) (WrappedTemplate, error) {
	t, err := w.Template.Parse(text)
	return WrappedTemplate{t}, err
}

func (w WrappedTemplate) Execute(wr io.Writer, data any) error {
	// HACK: Allows for nested template expansions that are not usually
	// possible.
	CurrentTemplate = w.Template
	err := w.Template.Execute(wr, data)
	CurrentTemplate = nil

	return err
}

func GetTemplate() WrappedTemplate {
	return WrappedTemplate{template.Must(baseTemplate.Clone())}
}

func ExpandString(name, content string, data map[string]any) (string, error) {
	if len(content) == 0 {
		return content, nil
	}

	t, err := template.Must(baseTemplate.Clone()).
		Parse(content)

	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	err = t.Execute(&buffer, data)

	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
