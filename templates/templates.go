package templates

import (
	"bytes"
	"fmt"
	"github.com/veigaribo/template"
)

var Funcs = template.FuncMap{
	"add": TemplateAdd,
	"sub": TemplateSub,
	"mul": TemplateMul,
	"div": TemplateDiv,
	"rem": TemplateRem,

	"join": TemplateJoinFn,

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

	"ismap": TemplateIsMap,
	"isstr": TemplateIsStr,
	"isint": TemplateIsInt,
	"isarr": TemplateIsArr,

	"map":    TemplateMap,
	"list":   TemplateList,
	"set":    TemplateSet,
	"append": TemplateAppend,
	"slice":  TemplateSlice,

	"jq1": TemplateJq1,
	"jqn": TemplateJqN,

	"err":   TemplateErr,
	"dump":  TemplateDump,
	"probe": TemplateProbe,

	"toml": TemplateToToml,
	"yaml": TemplateToYaml,
	"json": TemplateToJson,
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
		TemplateJoinT,
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

// Converts containers to pointers to containers.
func PrepareData(data any) any {
	switch val := data.(type) {
	case []any:
		newSlice := make([]any, 0, len(val))

		for _, item := range val {
			newSlice = append(newSlice, PrepareData(item))
		}

		return &newSlice
	case map[string]any:
		newMap := make(map[string]any)

		for key, value := range val {
			newMap[key] = PrepareData(value)
		}

		return &newMap
	default:
		return val
	}
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

// The template works with pointers to containers. Often we need to
// convert to embedded containers. Should do the opposite of
// `PrepareData`.
func resolvePointers(data any) any {
	switch val := data.(type) {
	case *[]any:
		newSlice := make([]any, 0, len(*val))

		for _, item := range *val {
			newSlice = append(newSlice, resolvePointers(item))
		}

		return newSlice
	case *map[string]any:
		newMap := make(map[string]any)

		for key, value := range *val {
			newMap[key] = resolvePointers(value)
		}

		return newMap
	default:
		return val
	}
}
