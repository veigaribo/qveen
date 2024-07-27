package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"text/template"

	"github.com/veigaribo/qveen/params"
	"github.com/veigaribo/qveen/prompts"
	"github.com/veigaribo/qveen/templates"
	"github.com/veigaribo/qveen/utils"
)

type Render1Options struct {
	ParamsPath   string
	TemplatePath string
	OutputPath   string
	MetaKey      string
}

func Render1(opts Render1Options) {
	paramsReader, err := utils.OpenFileOrUrl(opts.ParamsPath)

	if err != nil {
		panic(fmt.Errorf("Failed to open parameter file: %w", err))
	}

	p, err := params.ParseParams(paramsReader, params.ParseParamsOptions{
		MetaKey: opts.MetaKey,
	})

	if err != nil {
		panic(fmt.Errorf("Failed to parse parameter file: %w", err))
	}

	render(renderOptions{
		Params:           &p,
		TemplatePathFlag: opts.TemplatePath,
		OutputPathFlag:   opts.OutputPath,
		MetaKey:          opts.MetaKey,
	})
}

type RenderNOptions struct {
	ParamsPaths   []string
	OutputDirPath string
	MetaKey       string
}

func RenderN(opts RenderNOptions) {
	if !IsPrefix(opts.OutputDirPath) {
		panic(fmt.Errorf("Output path in multi-file mode can only be a prefix, but received '%s'", opts.OutputDirPath))
	}

	// Store stdin so it can be used multiple times.

	var stdinParams *params.ParsingParams
	if slices.Contains(opts.ParamsPaths, "-") {
		stdinParamsValue, err := params.ParseParams(
			os.Stdin,
			params.ParseParamsOptions{MetaKey: opts.MetaKey},
		)

		if err != nil {
			panic(fmt.Errorf("Failed to parse stdin: %w", err))
		}

		stdinParams = &stdinParamsValue
	}

	for i, paramsPath := range opts.ParamsPaths {
		var p params.ParsingParams

		if utils.IsStd(paramsPath) {
			p = *stdinParams
		} else {
			paramsReader, err := utils.OpenFileOrUrl(paramsPath)

			if err != nil {
				panic(fmt.Errorf("Failed to open parameter file '%s' (#%d): %w", paramsPath, i, err))
			}

			p, err = params.ParseParams(
				paramsReader,
				params.ParseParamsOptions{MetaKey: opts.MetaKey},
			)

			if err != nil {
				panic(fmt.Errorf("Failed to parse parameter file '%s' (#%d): %w", paramsPath, i, err))
			}
		}

		render(renderOptions{
			Params:           &p,
			TemplatePathFlag: "",
			OutputPathFlag:   "",
			MetaKey:          opts.MetaKey,
		})
	}
}

type renderOptions struct {
	Params           *params.ParsingParams
	TemplatePathFlag string
	OutputPathFlag   string
	MetaKey          string
}

// Ad-hoc function to do what both the above methods have in common.
func render(opts renderOptions) {
	opts.Params.ExpandPromptParams(opts.MetaKey)
	err := doPrompt(opts.Params.Prompt, opts.Params.Data)

	if err != nil {
		panic(fmt.Errorf("Failed to run prompts: %w", err))
	}

	err = opts.Params.ExpandParams(opts.MetaKey)

	templatePathFlag, err := templates.ExpandString(
		"--template",
		opts.TemplatePathFlag,
		opts.Params.Data,
	)

	if err != nil {
		panic(fmt.Errorf("Failed to expand template path: %w", err))
	}

	outputPathFlag, err := templates.ExpandString(
		"--output",
		opts.OutputPathFlag,
		opts.Params.Data,
	)

	if err != nil {
		panic(fmt.Errorf("Failed to expand output path: %w", err))
	}

	var templatePathFile = opts.Params.Template

	templatePath := utils.FirstOf(templatePathFlag, templatePathFile)

	if templatePath == "" {
		panic("Missing template file path.")
	}

	templateReader, err := utils.OpenFileOrUrl(templatePath)

	if err != nil {
		panic(fmt.Errorf("Failed to open template file: %w", err))
	}

	templateData, err := io.ReadAll(templateReader)

	if err != nil {
		panic(fmt.Errorf("Failed to read template file: %w", err))
	}

	t := template.Must(
		template.
			New(templatePath).
			Delims(templates.LeftDelim, templates.RightDelim).
			Funcs(templates.Funcs).
			Parse(string(templateData)),
	)

	var outputPathFile = opts.Params.Output

	var outputLoc OutputLocation
	outputLoc.Add(outputPathFile)
	outputLoc.Add(outputPathFlag)

	outputPath, err := outputLoc.Path()

	if err != nil {
		panic(fmt.Errorf("Failed to generate output path: %w", err))
	}

	output, err := utils.FileWriter(outputPath)

	if err != nil {
		panic(fmt.Errorf("Failed to create output file: %w", err))
	}

	err = t.Execute(output, opts.Params.Data)

	if err != nil {
		panic(fmt.Errorf("Failed to execute template: %w", err))
	}

	if utils.IsStd(outputPath) {
		// A little spacing for us humans.
		fmt.Fprintln(os.Stderr, "")
	}
	fmt.Fprintln(os.Stderr, templatePath, "->", outputPath)
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
