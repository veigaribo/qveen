package templates

import (
	"bytes"
	"fmt"
	"github.com/veigaribo/template"
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

	"escapebackslash": EscapeBackslash,
	"escapedouble":    EscapeDouble,
	"escapehtml":      EscapeHtml,
	"repl":            Replace,

	"map":  TemplateMap,
	"list": TemplateList,

	"jq1": TemplateJq1,
	"jqn": TemplateJqN,

	"err": TemplateErr,
}

var LeftDelim string = ""
var RightDelim string = ""

var baseTemplate *template.Template

func Init() {
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
