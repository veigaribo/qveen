package templates

import (
	"fmt"
)

func TemplateMap(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("Tried to construct map with an odd number of arguments (%d).", len(values))
	}

	m := make(map[string]any)

	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		value := values[i+1]

		m[key] = value
	}

	return m, nil
}

func TemplateList(values ...any) ([]any, error) {
	s := make([]any, 0, len(values))

	for _, value := range values {
		s = append(s, value)
	}

	return s, nil
}
