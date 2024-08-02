package templates

import (
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// Makes it easier to generate a parameter file from another parameter
// file.
func TemplateToToml(obj any) (string, error) {
	var builder strings.Builder

	encoder := toml.NewEncoder(&builder)
	encoder.SetIndentTables(false)

	err := encoder.Encode(resolvePointers(obj))

	if err != nil {
		return "", err
	}

	result := strings.TrimSpace(builder.String())
	return result, nil
}
