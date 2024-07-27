package params

import (
	"fmt"
	"io"
	"slices"

	"github.com/pelletier/go-toml/v2"
	"github.com/veigaribo/qveen/prompts"
)

type ParsingParams struct {
	Template string
	Output   string
	Prompt   []prompts.Prompt
	Data     map[string]any
}

type ParseParamsOptions struct {
	MetaKey string
}

func ParseParams(
	input io.Reader,
	opts ParseParamsOptions,
) (ParsingParams, error) {
	if opts.MetaKey == "" {
		opts.MetaKey = "meta"
	}

	var params ParsingParams
	var err error

	err = params.ParseGeneralParams(input)

	if err != nil {
		return params, err
	}

	err = params.ParseMetaParams(opts)

	if err != nil {
		return params, err
	}

	return params, nil
}

func (params *ParsingParams) ParseGeneralParams(
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

func (params *ParsingParams) ParseMetaParams(opts ParseParamsOptions) error {
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

	templateRaw, ok := meta["template"]

	if ok {
		template, ok := templateRaw.(string)

		if !ok {
			return MakeParamError(
				[]any{opts.MetaKey, "template"},
				"field present but does not contain a string.",
			)
		}

		params.Template = template
	}

	outputRaw, ok := meta["output"]

	if ok {
		output, ok := outputRaw.(string)

		if !ok {
			return MakeParamError(
				[]any{opts.MetaKey, "output"},
				"field present but does not contain a string.",
			)
		}

		params.Output = output
	}

	promptRaw, ok := meta["prompt"]

	if ok {
		prompt, ok := promptRaw.([]any)

		if !ok {
			return MakeParamError(
				[]any{opts.MetaKey, "prompt"},
				"field present but does not contain an array.",
			)
		}

		for i, entryRaw := range prompt {
			entry, ok := entryRaw.(map[string]any)

			if !ok {
				return MakeParamError(
					[]any{opts.MetaKey, "prompt", i},
					"field present but does not contain a map.",
				)
			}

			prompt, err := parseMetaPrompt(entry,
				[]any{opts.MetaKey, "prompt", i},
			)

			if err != nil {
				return err
			}

			params.Prompt = append(params.Prompt, prompt)
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
