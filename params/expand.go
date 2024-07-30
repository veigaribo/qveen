package params

import (
	"github.com/veigaribo/qveen/templates"
	"github.com/veigaribo/qveen/utils"
)

// Prompt param expansion must be done earlier to display the correct
// prompts...
func (p *Params) ExpandPromptParams(metaKey string) error {
	var err error

	for i := range p.Prompt {
		entry := &p.Prompt[i]

		templateName := func(field string) string {
			return utils.PathString(
				[]any{metaKey, "prompts", i, field},
			)
		}

		entry.Name, err = templates.ExpandString(
			templateName("name"),
			entry.Name,
			p.Data,
		)

		if err != nil {
			return err
		}

		entry.Title, err = templates.ExpandString(
			templateName("title"),
			entry.Title,
			p.Data,
		)

		if err != nil {
			return err
		}

		// `kind` intentionally left as is.
	}

	return nil
}

// Because maps are not addressable in Go, we need to keep a reference
// to the map + the key in order to update it. In our case, we may
// also need to patch a slice, so that's even more fun.
type ContainerPtr struct {
	Container any
	Key       any
	Path      []any
}

func MakeContainerPtr(container any, key any, path []any) ContainerPtr {
	return ContainerPtr{
		Container: container,
		Key:       key,
		Path:      path,
	}
}

func (c ContainerPtr) Set(value any) {
	if m, ok := c.Container.(map[string]any); ok {
		key := c.Key.(string)

		m[key] = value
	} else {
		arr := c.Container.([]any)
		key := c.Key.(int)

		arr[key] = value
	}
}

// Recursively expands strings.
func expandParamsVisit(
	data map[string]any,
	ptr ContainerPtr,
	value any,
) error {
	if str, ok := value.(string); ok {
		expanded, err := templates.ExpandString(
			utils.PathString(append(ptr.Path, ptr.Key)),
			str,
			data,
		)

		if err != nil {
			return err
		}

		ptr.Set(expanded)
	} else if m, ok := value.(map[string]any); ok {
		for k, v := range m {
			ptr := MakeContainerPtr(m, k, append(ptr.Path, k))
			return expandParamsVisit(data, ptr, v)
		}
	} else if s, ok := value.([]any); ok {
		for i, v := range s {
			ptr := MakeContainerPtr(s, i, append(ptr.Path, i))
			return expandParamsVisit(data, ptr, v)
		}
	}

	return nil
}

// ...other params must be expanded later to use the values of the
// prompts.
func (p *Params) ExpandParams(metaKey string) error {
	var err error

	// General fields.

	for k, v := range p.Data {
		ptr := MakeContainerPtr(p.Data, k, []any{})
		expandParamsVisit(p.Data, ptr, v)
	}

	// Meta fields.

	for i := range p.Pairs {
		pair := &p.Pairs[i]

		metaTemplateName := func(field string) string {
			return utils.PathString(append(pair.Path, field))
		}

		pair.Template.Path, err = templates.ExpandString(
			metaTemplateName("template"),
			pair.Template.Path,
			p.Data,
		)

		if err != nil {
			return err
		}

		pair.Output.Path, err = templates.ExpandString(
			metaTemplateName("output"),
			pair.Output.Path,
			p.Data,
		)

		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	return nil
}
