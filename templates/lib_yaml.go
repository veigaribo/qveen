package templates

import (
	"encoding/json"
	"strings"

	"gopkg.in/yaml.v3"
)

// Makes it easier to generate a parameter file from another parameter
// file.
func TemplateToYaml(obj any) (string, error) {
	var builder strings.Builder

	encoder := yaml.NewEncoder(&builder)
	err := encoder.Encode(obj)

	if err != nil {
		return "", err
	}

	result := strings.TrimSpace(builder.String())
	return result, nil
}

// Makes it easier to generate a parameter file from another parameter
// file.
func TemplateToJson(obj any) (string, error) {
	var builder strings.Builder

	encoder := json.NewEncoder(&builder)
	err := encoder.Encode(obj)

	if err != nil {
		return "", err
	}

	result := strings.TrimSpace(builder.String())
	return result, nil
}
