package templates

import "github.com/itchyny/gojq"

// Runs jq query and returns the first result.
func TemplateJq1(query string, obj any) (any, error) {
	q, err := gojq.Parse(query)

	if err != nil {
		return nil, err
	}

	iter := q.Run(obj)

	v, ok := iter.Next()
	if !ok {
		return nil, nil
	}
	if err, ok := v.(error); ok {
		if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
			return nil, nil
		}

		return "", err
	}

	return v, nil
}

// Runs jq query and returns all results.
func TemplateJqN(query string, obj any) ([]any, error) {
	q, err := gojq.Parse(query)
	results := make([]any, 0)

	if err != nil {
		return results, err
	}

	iter := q.Run(obj)

	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
				break
			}

			return results, err
		}

		results = append(results, v)
	}

	return results, nil
}
