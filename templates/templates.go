package templates

import (
	"bytes"
	"fmt"
	"text/template"
)

var Funcs = template.FuncMap{
	"uppercase":    UpperCase,
	"lowercase":    LowerCase,
	"titlecase":    TitleCase,
	"pascalcase":   PascalCase,
	"camelcase":    CamelCase,
	"snakecase":    SnakeCase,
	"kebabcase":    KebabCase,
	"constcase":    ConstantCase,
	"dotcase":      DotCase,
	"sentencecase": SentenceCase,

	"map": Map,
}

var LeftDelim string = ""
var RightDelim string = ""

var baseTemplate *template.Template

func init() {
	var err error

	baseTemplate = template.
		New("qveen_base").
		Delims(LeftDelim, RightDelim).
		Funcs(Funcs)

	builtinTemplates := []string{
		TemplateJoin,
	}

	for i, builtin := range builtinTemplates {
		baseTemplate, err = baseTemplate.Parse(builtin)

		if err != nil {
			panic(fmt.Errorf("Failed to parse builtin template #%d!", i))
		}
	}
}

func GetTemplate() *template.Template {
	return template.Must(baseTemplate.Clone())
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
