package params

import (
	"fmt"
	"io"
	"path"
	"slices"

	"github.com/pelletier/go-toml/v2"
	"github.com/veigaribo/qveen/prompts"
	"gopkg.in/yaml.v3"
)

type ParamsPathFrom = uint

const (
	ParamsPathFromCwd ParamsPathFrom = iota
	ParamsPathFromParams
)

type ParamsPath struct {
	Path string
	From ParamsPathFrom
}

func (pp ParamsPath) IsEmpty() bool {
	return pp.Path == ""
}

func (pp ParamsPath) Resolve(pathToParams string) string {
	if pp.IsEmpty() {
		return ""
	}

	if pp.From == ParamsPathFromCwd {
		return pp.Path
	}

	// Must be ParamsPathFromParams.
	paramsDir := path.Dir(pathToParams)
	return path.Join(paramsDir, pp.Path)
}

type ParamsPair struct {
	Template ParamsPath
	Output   ParamsPath

	// Because a pair may be found in multiple keys (`meta` or
	// `meta.pairs[#]`), store it here so we can show appropriate
	// error messages after parsing.
	Path []any
}

type Params struct {
	Data   map[string]any
	Pairs  []ParamsPair
	Prompt []prompts.Prompt

	TemplateLeftDelim  string
	TemplateRightDelim string
	TemplateCase       string
}

type ParamsFormat string

const (
	ParamsTomlFormat ParamsFormat = "toml"
	ParamsYamlFormat              = "yaml/json"
)

type ParseParamsOptions struct {
	MetaKey string
}

func ParseParams(
	input io.Reader,
	format ParamsFormat,
	opts ParseParamsOptions,
) (Params, error) {
	if opts.MetaKey == "" {
		opts.MetaKey = "meta"
	}

	var params Params
	var err error

	err = params.ParseGeneral(input, format)

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
	format ParamsFormat,
) error {
	bytes, err := io.ReadAll(input)

	if err != nil {
		return err
	}

	switch format {
	case ParamsTomlFormat:
		err = toml.Unmarshal(bytes, &params.Data)
	case ParamsYamlFormat:
		err = yaml.Unmarshal(bytes, &params.Data)
	default:
		panic(fmt.Errorf("Unrecognized format '%q'", format))
	}

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
		return MakeMetaWrongTypeError([]any{opts.MetaKey})
	}

	err := params.parseMetaPairs(meta, []any{opts.MetaKey})

	if err != nil {
		return err
	}

	err = params.parseMetaPrompts(meta, []any{opts.MetaKey})

	if err != nil {
		return err
	}

	leftDelimRaw, ok := meta["left_delim"]

	if ok {
		leftDelim, ok := leftDelimRaw.(string)

		if !ok {
			return MakeMetaLeftDelimWrongTypeError([]any{opts.MetaKey, "left_delim"})
		}

		params.TemplateLeftDelim = leftDelim
	}

	rightDelimRaw, ok := meta["right_delim"]

	if ok {
		rightDelim, ok := rightDelimRaw.(string)

		if !ok {
			return MakeMetaRightDelimWrongTypeError([]any{opts.MetaKey, "right_delim"})
		}

		params.TemplateRightDelim = rightDelim
	}

	caseRaw, ok := meta["case"]

	if ok {
		kase, ok := caseRaw.(string)

		if !ok {
			return MakeMetaCaseWrongTypeError([]any{opts.MetaKey, "case"})
		}

		params.TemplateCase = kase
	}

	return nil
}

// We need to explain the concept of variance to the compiler
// apparently.
func rerr[T error](f func([]any) T) func([]any) error {
	return func(xs []any) error {
		return f(xs)
	}
}

func (p *Params) parseMetaPairs(
	meta map[string]any, path []any,
) error {
	var rootTemplate, rootOutput ParamsPath

	templateRaw, ok := meta["template"]

	if ok {
		template, err := parsePath(
			templateRaw,
			append(path, "template"),
			mkParsePathErrors{
				WrongType:          rerr(MakeMetaRootTemplateWrongTypeError),
				TablePathMissing:   rerr(MakeMetaRootTemplatePathMissingError),
				TablePathWrongType: rerr(MakeMetaRootTemplatePathWrongTypeError),
				TableFromWrongType: rerr(MakeMetaRootTemplateFromWrongTypeError),
				TableFromInvalid:   rerr(MakeMetaRootTemplateFromInvalidError),
			},
		)

		if err != nil {
			return err
		}

		rootTemplate = template
	}

	outputRaw, ok := meta["output"]

	if ok {
		output, err := parsePath(
			outputRaw,
			append(path, "output"),
			mkParsePathErrors{
				WrongType:          rerr(MakeMetaRootOutputWrongTypeError),
				TablePathMissing:   rerr(MakeMetaRootOutputPathMissingError),
				TablePathWrongType: rerr(MakeMetaRootOutputPathWrongTypeError),
				TableFromWrongType: rerr(MakeMetaRootOutputFromWrongTypeError),
				TableFromInvalid:   rerr(MakeMetaRootOutputFromInvalidError),
			},
		)

		if err != nil {
			return err
		}

		rootOutput = output
	}

	if !rootTemplate.IsEmpty() || !rootOutput.IsEmpty() {
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
			return MakeMetaPairsWrongTypeError(append(path, "pairs"))
		}

		for i, pairRaw := range pairs {
			entry, ok := pairRaw.(map[string]any)

			if !ok {
				return MakeMetaPairWrongTypeError(append(path, "pairs", i))
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

		if first.Template.IsEmpty() {
			return MakeMetaRootTemplateMissingInMultipleError(append(first.Path, "template"))
		}

		if first.Output.IsEmpty() {
			return MakeMetaRootOutputMissingInMultipleError(append(first.Path, "output"))
		}
	}

	return nil
}

func parseMetaPair(
	entry map[string]any, path []any,
) (ParamsPair, error) {
	var pair ParamsPair
	var err error

	templateRaw, ok := entry["template"]

	if !ok {
		return pair, MakeMetaPairTemplateMissingError(append(path, "template"))
	}

	pair.Template, err = parsePath(
		templateRaw,
		append(path, "template"),
		mkParsePathErrors{
			WrongType:          rerr(MakeMetaPairTemplateWrongTypeError),
			TablePathMissing:   rerr(MakeMetaPairTemplatePathMissingError),
			TablePathWrongType: rerr(MakeMetaPairTemplatePathWrongTypeError),
			TableFromWrongType: rerr(MakeMetaPairTemplateFromWrongTypeError),
			TableFromInvalid:   rerr(MakeMetaPairTemplateFromInvalidError),
		},
	)

	if err != nil {
		return pair, err
	}

	outputRaw, ok := entry["output"]

	if !ok {
		return pair, MakeMetaPairOutputMissingError(append(path, "output"))
	}

	pair.Output, err = parsePath(
		outputRaw,
		append(path, "output"),
		mkParsePathErrors{
			WrongType:          rerr(MakeMetaPairOutputWrongTypeError),
			TablePathMissing:   rerr(MakeMetaPairOutputPathMissingError),
			TablePathWrongType: rerr(MakeMetaPairOutputPathWrongTypeError),
			TableFromWrongType: rerr(MakeMetaPairOutputFromWrongTypeError),
			TableFromInvalid:   rerr(MakeMetaPairOutputFromInvalidError),
		},
	)

	pair.Path = path
	return pair, nil
}

func (p *Params) parseMetaPrompts(
	meta map[string]any, path []any,
) error {
	promptRaw, ok := meta["prompts"]

	if ok {
		prompt, ok := promptRaw.([]any)

		if !ok {
			return MakeMetaPromptsWrongTypeError(append(path, "prompts"))
		}

		for i, entryRaw := range prompt {
			entry, ok := entryRaw.(map[string]any)

			if !ok {
				return MakeMetaPromptWrongTypeError(append(path, "prompts", i))
			}

			prompt, err := parseMetaPrompt(entry,
				append(path, "prompts", i),
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
		return prompt, MakeMetaPromptNameMissingError(append(path, "name"))
	}

	name, ok := nameRaw.(string)

	if !ok {
		return prompt, MakeMetaPromptNameWrongTypeError(append(path, "name"))
	}

	kindRaw, ok := entry["kind"]

	if !ok {
		goto postKind
	}

	kind, ok = kindRaw.(string)

	if !ok {
		return prompt, MakeMetaPromptKindWrongTypeError(append(path, "kind"))
	}

	if !slices.Contains(prompts.SupportedPromptKinds, kind) {
		return prompt, MakeMetaPromptKindInvalidError(append(path, "kind"))
	}

postKind:

	titleRaw, ok := entry["title"]

	if !ok {
		goto postTitle
	}

	title, ok = titleRaw.(string)

	if !ok {
		return prompt, MakeMetaPromptTitleWrongTypeError(append(path, "title"))
	}

postTitle:

	var specific any = nil

	switch kind {
	case "select":
		optionsRaw, ok := entry["options"]

		if !ok {
			return prompt, MakeMetaPromptOptionsMissingError(append(path, "options"))
		}

		options, ok := optionsRaw.([]any)

		if !ok {
			return prompt, MakeMetaPromptOptionsWrongTypeError(append(path, "options"))
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
					return prompt, MakeMetaPromptOptionTitleMissingError(append(path, "options", i, "title"))
				}

				title, ok = titleRaw.(string)

				if !ok {
					return prompt, MakeMetaPromptOptionTitleWrongTypeError(append(path, "options", i, "title"))
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
			} else {
				return prompt, MakeMetaPromptWrongTypeError(append(path, "options", i))
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

type mkParsePathErrors struct {
	WrongType          func(path []any) error
	TablePathMissing   func(path []any) error
	TablePathWrongType func(path []any) error
	TableFromWrongType func(path []any) error
	TableFromInvalid   func(path []any) error
}

func parsePath(
	obj any,
	path []any,
	mkerr mkParsePathErrors,
) (ParamsPath, error) {

	result := ParamsPath{
		From: ParamsPathFromCwd,
	}

	var ok bool

	result.Path, ok = obj.(string)

	if ok {
		return result, nil
	}

	m, ok := obj.(map[string]any)

	if !ok {
		return result, mkerr.WrongType(path)
	}

	pathRaw, ok := m["path"]

	if !ok {
		return result, mkerr.TablePathMissing(append(path, "path"))
	}

	result.Path, ok = pathRaw.(string)

	if !ok {
		return result, mkerr.TablePathWrongType(append(path, "path"))
	}

	fromRaw, ok := m["from"]
	var from string

	if !ok {
		goto postFrom
	}

	from, ok = fromRaw.(string)

	if !ok {
		return result, mkerr.TableFromWrongType(append(path, "from"))
	}

	switch from {
	case "params":
		result.From = ParamsPathFromParams
	case "cwd":
		result.From = ParamsPathFromCwd
	default:
		return result, mkerr.TableFromInvalid(append(path, "from"))
	}

postFrom:
	return result, nil
}
