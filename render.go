package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"unicode"

	"github.com/veigaribo/qveen/params"
	"github.com/veigaribo/qveen/prompts"
	"github.com/veigaribo/qveen/templates"
	"github.com/veigaribo/qveen/utils"
)

type RenderOptions struct {
	ParamsPath   string
	ParamsFormat string
	TemplatePath string
	OutputPath   string
	MetaKey      string
	PromptValues map[string]string
	Overwrite    bool

	TemplateLeftDelim  string
	TemplateRightDelim string
	TemplateCase       string
}

func Render(opts RenderOptions) {
	paramsReader, err := utils.OpenFileOrUrl(opts.ParamsPath)

	if err != nil {
		panic(fmt.Errorf("Failed to open parameter file: %w", err))
	}

	var paramsFormat params.ParamsFormat

	switch opts.ParamsFormat {
	case "toml":
		paramsFormat = params.ParamsTomlFormat
	case "yaml":
	case "json":
		paramsFormat = params.ParamsYamlFormat
	case "":
		maybeParamsFormat := params.GuessFormat(opts.ParamsPath)

		if maybeParamsFormat == nil {
			panic(fmt.Errorf("Could not guess params format from file name: %s", opts.ParamsPath))
		}

		paramsFormat = *maybeParamsFormat
	default:
		panic(fmt.Errorf("Unrecognized format '%s'", opts.ParamsFormat))
	}

	p, err := params.ParseParams(
		paramsReader,
		paramsFormat,

		params.ParseParamsOptions{
			MetaKey: opts.MetaKey,
		},
	)

	if err != nil {
		panic(fmt.Errorf("Failed to parse parameter file: %w", err))
	}

	if len(p.Pairs) == 0 {
		// Nothing to do.
		fmt.Fprintf(os.Stderr, "Nothing to do.\n")
		return
	}

	templateLeftDelim := utils.FirstOf(
		opts.TemplateLeftDelim,
		p.TemplateLeftDelim,
	)

	templates.LeftDelim = templateLeftDelim

	templateRightDelim := utils.FirstOf(
		opts.TemplateRightDelim,
		p.TemplateRightDelim,
	)

	templates.RightDelim = templateRightDelim

	templateCase := utils.FirstOf(opts.TemplateCase, p.TemplateCase)

	switch templateCase {
	case "turkish":
		templates.Case = unicode.TurkishCase
	case "azeri":
		templates.Case = unicode.AzeriCase
	case "":
		break
	default:
		panic(fmt.Errorf("Invalid value for `case`: '%s'. Expected 'turkish' or 'azeri' (or empty).", templateCase))
	}

	templates.Init()
	p.ExpandPromptParams(opts.MetaKey)

	for i := range p.Prompt {
		prompt := &p.Prompt[i]
		prefill, ok := opts.PromptValues[prompt.Name]

		if ok {
			err := prompt.TryPrefill(prefill)

			if err != nil {
				panic(fmt.Errorf("Failed to prefill prompt '%s' with '%s': %w", prompt.Name, prefill, err))
			}
		}
	}

	err = doPrompt(p.Prompt, p.Data)

	if err != nil {
		panic(fmt.Errorf("Failed to run prompts: %w", err))
	}

	err = p.ExpandParams(opts.MetaKey)

	isSinglePair := len(p.Pairs) == 1

	var templatePathFlag, outputPathFlag string

	if isSinglePair {
		if opts.TemplatePath != "" {
			templatePathFlag, err = templates.ExpandString(
				"--template",
				opts.TemplatePath,
				p.Data,
			)

			if err != nil {
				panic(fmt.Errorf("Failed to expand template path: %w", err))
			}
		}

		if opts.OutputPath != "" {
			outputPathFlag, err = templates.ExpandString(
				"--output",
				opts.OutputPath,
				p.Data,
			)

			if err != nil {
				panic(fmt.Errorf("Failed to expand output path: %w", err))
			}
		}
	} else {
		if opts.TemplatePath != "" {
			fmt.Fprintf(os.Stderr, "Ignoring template flag for multiple pairs.")
		}

		if opts.OutputPath != "" {
			if utils.IsExplicitDir(opts.OutputPath) {
				outputPathFlag = opts.OutputPath
			} else {
				fmt.Fprintf(os.Stderr, "Ignoring non-prefix output flag for multiple pairs.")
			}
		}
	}

	for i, pair := range p.Pairs {
		templatePathParams := pair.Template
		templatePath := utils.FirstOf(
			templatePathFlag,
			templatePathParams.Resolve(opts.ParamsPath),
		)

		if templatePath == "" {
			if isSinglePair {
				panic(errors.New("Missing template file path."))
			} else {
				panic(fmt.Errorf("Missing template file path for pair #%d.", i))
			}
		}

		templateReader, err := utils.OpenFileOrUrl(templatePath)

		if err != nil {
			if isSinglePair {
				panic(fmt.Errorf("Failed to open template file: %w", err))
			} else {
				panic(fmt.Errorf("Failed to open template file for pair #%d: %w", i, err))
			}
		}

		templateData, err := io.ReadAll(templateReader)

		if err != nil {
			if isSinglePair {
				panic(fmt.Errorf("Failed to read template file: %w", err))
			} else {
				panic(fmt.Errorf("Failed to read template file for pair #%d: %w", i, err))
			}
		}

		t, err := templates.GetTemplate().Parse(string(templateData))

		if err != nil {
			panic(fmt.Errorf("Failed to parse template: %w", err))
		}

		var outputPathParams = pair.Output

		var outputLoc OutputLocation
		outputLoc.Add(outputPathParams.Resolve(opts.ParamsPath))
		outputLoc.Add(outputPathFlag)

		outputPath, err := outputLoc.Path()

		if err != nil {
			if isSinglePair {
				panic(fmt.Errorf("Failed to generate output path: %w", err))
			} else {
				panic(fmt.Errorf("Failed to generate output path for pair #%d: %w", i, err))
			}
		}

		output, err := utils.FileWriter(outputPath, opts.Overwrite)

		if err != nil {
			if isSinglePair {
				panic(fmt.Errorf("Failed to create output file: %w", err))
			} else {
				panic(fmt.Errorf("Failed to create output file for pair #%d: %w", i, err))
			}
		}

		err = t.Execute(output, templates.PrepareData(p.Data))

		if err != nil {
			if isSinglePair {
				panic(fmt.Errorf("Failed to execute template: %w", err))
			} else {
				panic(fmt.Errorf("Failed to execute template for pair #%d: %w", i, err))
			}
		}

		fmt.Fprintln(os.Stderr, i, templatePath, "->", outputPath)
	}
}

func doPrompt(ps []prompts.Prompt, out map[string]any) error {
	prompted, err := prompts.DoPrompt(ps)

	if err != nil {
		return err
	}

	for key, value := range prompted {
		out[key] = value
	}

	return nil
}
