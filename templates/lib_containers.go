package templates

import (
	"fmt"
)

func TemplateMap(values ...any) (*map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("Tried to construct map with an odd number of arguments (%d).", len(values))
	}

	m := make(map[string]any)

	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		value := values[i+1]

		m[key] = value
	}

	return &m, nil
}

func TemplateList(values ...any) (*[]any, error) {
	s := make([]any, 0, len(values))

	for _, value := range values {
		s = append(s, value)
	}

	return &s, nil
}

func TemplateSet(container any, key any, value any) (string, error) {
	switch c := container.(type) {
	case *[]any:
		ikey, ok := key.(int)

		if !ok {
			return "", fmt.Errorf("Tried to use key of type '%[1]T' (%[1]q) to index a slice. Expected an int.", key)
		}

		if ikey < len(*c) {
			newSlice := make([]any, ikey)
			copy(newSlice, *c)
			*c = newSlice
		}

		(*c)[ikey] = value
	case *map[string]any:
		skey, ok := key.(string)

		if !ok {
			return "", fmt.Errorf("Tried to use key of type '%[1]T' (%[1]q) to index a map. Expected a string.", key)
		}

		(*c)[skey] = value
	default:
		return "", fmt.Errorf("Container '%[1]q' of type %[1]T is neither a list nor a map. Cannot set.", container)
	}

	return "", nil
}

func TemplateAppend(s *[]any, value any) (string, error) {
	*s = append(*s, value)
	return "", nil
}

func TemplateSlice(s *[]any, idx ...int) (*[]any, error) {
	if len(idx) < 1 || len(idx) > 2 {
		return nil, fmt.Errorf("`slice` expected between 1 and 2 indexes. Received %d (%v)", len(idx), idx)
	}

	switch len(idx) {
	case 1:
		slice := (*s)[idx[0]:]
		return &slice, nil
	case 2:
		slice := (*s)[idx[0]:idx[1]]
		return &slice, nil
	}

	return nil, nil
}
