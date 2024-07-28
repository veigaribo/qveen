package params

import (
	"fmt"
	"io"
	"slices"

	"github.com/pelletier/go-toml/v2"
	"github.com/veigaribo/qveen/prompts"
)

type ParamsPair struct {
	Template string
	Output   string

	// Because a pair may be found in multiple keys (`meta` or
	// `meta.pairs[#]`), store it here so we can show appropriate
	// error messages after parsing.
	Path []any
}

type Params struct {
	Data   map[string]any
	Pairs  []ParamsPair
	Prompt []prompts.Prompt
}

type ParseParamsOptions struct {
	MetaKey string
}

func ParseParams(
	input io.Reader,
	opts ParseParamsOptions,
) (Params, error) {
	if opts.MetaKey == "" {
		opts.MetaKey = "meta"
	}

	var params Params
	var err error

	err = params.ParseGeneral(input)

	if err != nil {
		return params, err
	}

	err = params.ParseMeta(opts)

	if err != nil {
		return params, err
	}

	return params, nil
}

func (params *Params) ParseGeneral(
	input io.Reader,
) error {
	bytes, err := io.ReadAll(input)

	if err != nil {
		return err
	}

	err = toml.Unmarshal(bytes, &params.Data)

	if err != nil {
		return err
	}

	return nil
}

func (params *Params) ParseMeta(opts ParseParamsOptions) error {
	metaRaw, ok := params.Data[opts.MetaKey]

	if !ok {
		return nil
	}

	meta, ok := metaRaw.(map[string]any)

	if !ok {
		return MakeParamError(
			[]any{opts.MetaKey},
			"field present but does not contain a table.",
		)
	}

	err := params.parseMetaPairs(meta, []any{opts.MetaKey})

	if err != nil {
		return err
	}

	err = params.parseMetaPrompts(meta, []any{opts.MetaKey})

	if err != nil {
		return err
	}

	return nil
}

func (p *Params) parseMetaPairs(
	meta map[string]any, path []any,
) error {
	var rootTemplate, rootOutput string

	templateRaw, ok := meta["template"]

	if ok {
		template, ok := templateRaw.(string)

		if !ok {
			return MakeParamError(
				append(path, "template"),
				"field present but does not contain a string.",
			)
		}

		rootTemplate = template
	}

	outputRaw, ok := meta["output"]

	if ok {
		output, ok := outputRaw.(string)

		if !ok {
			return MakeParamError(
				append(path, "output"),
				"field present but does not contain a string.",
			)
		}

		rootOutput = output
	}

	if rootTemplate != "" || rootOutput != "" {
		p.Pairs = append(p.Pairs, ParamsPair{
			Template: rootTemplate,
			Output:   rootOutput,
			Path:     path,
		})
	}

	pairsRaw, ok := meta["pairs"]

	if ok {
		pairs, ok := pairsRaw.([]any)

		if !ok {
			return MakeParamError(
				append(path, "pairs"),
				"field present but does not contain an array.",
			)
		}

		for i, pairRaw := range pairs {
			entry, ok := pairRaw.(map[string]any)

			if !ok {
				return MakeParamError(
					append(path, "pairs", i),
					"field present but does not contain a map.",
				)
			}

			pair, err := parseMetaPair(entry,
				append(path, "pairs", i),
			)

			if err != nil {
				return err
			}

			p.Pairs = append(p.Pairs, pair)
		}
	}

	if len(p.Pairs) > 1 {
		first := p.Pairs[0]

		if first.Template == "" {
			return MakeParamError(
				append(first.Path, "template"),
				"required field for multiple pairs missing.",
			)
		}

		if first.Output == "" {
			return MakeParamError(
				append(first.Path, "output"),
				"required field for multiple pairs missing.",
			)
		}
	}

	return nil
}

func parseMetaPair(
	entry map[string]any, path []any,
) (ParamsPair, error) {
	var pair ParamsPair

	templateRaw, ok := entry["template"]

	if !ok {
		return pair, MakeParamError(
			append(path, "template"),
			"required field missing.",
		)
	}

	pair.Template, ok = templateRaw.(string)

	if !ok {
		return pair, MakeParamError(
			append(path, "template"),
			"field present but does not contain a string.",
		)
	}

	outputRaw, ok := entry["output"]

	if !ok {
		return pair, MakeParamError(
			append(path, "output"),
			"required field missing.",
		)
	}

	pair.Output, ok = outputRaw.(string)

	if !ok {
		return pair, MakeParamError(
			append(path, "output"),
			"field present but does not contain a string.",
		)
	}

	pair.Path = path
	return pair, nil
}

func (p *Params) parseMetaPrompts(
	meta map[string]any, path []any,
) error {
	promptRaw, ok := meta["prompt"]

	if ok {
		prompt, ok := promptRaw.([]any)

		if !ok {
			return MakeParamError(
				append(path, "prompt"),
				"field present but does not contain an array.",
			)
		}

		for i, entryRaw := range prompt {
			entry, ok := entryRaw.(map[string]any)

			if !ok {
				return MakeParamError(
					append(path, "prompt", i),
					"field present but does not contain a map.",
				)
			}

			prompt, err := parseMetaPrompt(entry,
				append(path, "prompt", i),
			)

			if err != nil {
				return err
			}

			p.Prompt = append(p.Prompt, prompt)
		}
	}

	return nil
}

func parseMetaPrompt(
	entry map[string]any, path []any,
) (prompts.Prompt, error) {
	var prompt prompts.Prompt

	// Defaults.
	kind := prompts.SupportedPromptKinds[0]
	title := ""

	nameRaw, ok := entry["name"]

	if !ok {
		return prompt, MakeParamError(
			append(path, "name"),
			"required field missing.",
		)
	}

	name, ok := nameRaw.(string)

	if !ok {
		return prompt, MakeParamError(
			append(path, "name"),
			"field present but does not contain a string.",
		)
	}

	kindRaw, ok := entry["kind"]

	if !ok {
		goto postKind
	}

	kind, ok = kindRaw.(string)

	if !ok {
		return prompt, MakeParamError(
			append(path, "kind"),
			"field present but does not contain a string.",
		)
	}

	if !slices.Contains(prompts.SupportedPromptKinds, kind) {
		return prompt, MakeParamError(
			append(path, "kind"),
			fmt.Sprintf(
				"field does not contain one of the allowed values: %v.",
				prompts.SupportedPromptKinds,
			),
		)
	}

postKind:

	titleRaw, ok := entry["title"]

	if !ok {
		goto postTitle
	}

	title, ok = titleRaw.(string)

	if !ok {
		return prompt, MakeParamError(
			append(path, "title"),
			"field present but does not contain a string.",
		)
	}

postTitle:

	var specific any = nil

	switch kind {
	case "select":
		optionsRaw, ok := entry["options"]

		if !ok {
			return prompt, MakeParamError(
				append(path, "options"),
				"required field for `select` missing.",
			)
		}

		options, ok := optionsRaw.([]any)

		if !ok {
			return prompt, MakeParamError(
				append(path, "options"),
				"field present but does not contain an array.",
			)
		}

		var optionsNormalized []prompts.PromptSelectOption

		for i, option := range options {
			if optionStr, ok := option.(string); ok {
				optionsNormalized = append(optionsNormalized,
					prompts.PromptSelectOption{
						Title: optionStr,
						Value: optionStr,
					},
				)
			} else if optionMap, ok := option.(map[string]any); ok {
				titleRaw, ok = optionMap["title"]
				var title string

				if !ok {
					return prompt, MakeParamError(
						append(path, "options", i, "title"),
						"required field missing.",
					)
				}

				title, ok = titleRaw.(string)

				if !ok {
					return prompt, MakeParamError(
						append(path, "options", i, "title"),
						"field present but does not contain a string.",
					)
				}

				value, ok := optionMap["value"]

				if !ok {
					value = title
				}

				optionsNormalized = append(optionsNormalized,
					prompts.PromptSelectOption{
						Title: title,
						Value: value,
					},
				)
			}
		}

		specific = prompts.PromptSelectSpecific{
			Options: optionsNormalized,
		}
	}

	prompt.Name = name
	prompt.Kind = kind
	prompt.Title = title
	prompt.Specific = specific
	return prompt, nil
}
