package templates

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
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
}

var LeftDelim string = ""
var RightDelim string = ""

func ExpandString(name, content string, data map[string]any) (string, error) {
	if len(content) == 0 {
		return content, nil
	}

	t, err := template.
		New(name).
		Delims(LeftDelim, RightDelim).
		Funcs(Funcs).
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

func GetTemplate(pathOrUrl string) (string, error) {
	var bytes []byte
	var err error

	if strings.HasPrefix(pathOrUrl, "http://") ||
		strings.HasPrefix(pathOrUrl, "https://") {
		resp, err := http.Get(pathOrUrl)

		if err != nil {
			return "", err
		}

		bytes, err = io.ReadAll(resp.Body)
	} else {
		bytes, err = os.ReadFile(pathOrUrl)
	}

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
