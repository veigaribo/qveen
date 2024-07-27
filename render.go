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
)

type Render1Options struct {
	ParamsPath   string
	TemplatePath string
	OutputPath   string
	MetaKey      string
}

func Render1(opts Render1Options) {
	var paramsIo io.Reader

	paramsPathFlag := opts.ParamsPath

	if paramsPathFlag == "" || paramsPathFlag == "-" {
		paramsIo = os.Stdin
	} else {
		var err error
		paramsIo, err = os.Open(paramsPathFlag)

		if err != nil {
			panic(fmt.Errorf("Failed to open parameter file: %w", err))
		}
	}

	p, err := params.ParseParams(paramsIo, params.ParseParamsOptions{
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

	var paramsIo io.Reader

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

		if paramsPath == "-" {
			p = *stdinParams
		} else {
			var err error
			paramsIo, err = os.Open(paramsPath)

			if err != nil {
				panic(fmt.Errorf("Failed to open parameter file '%s' (#%d): %w", paramsPath, i, err))
			}

			p, err = params.ParseParams(
				paramsIo,
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

	var templatePath string
	if templatePathFlag == "" {
		if templatePathFile == "" {
			panic("Missing template file path.")
		} else {
			templatePath = templatePathFile
		}
	} else {
		templatePath = templatePathFlag
	}

	templateData, err := templates.GetTemplate(templatePath)

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

	output, err := outputLoc.Writer()

	if err != nil {
		panic(fmt.Errorf("Failed to create output file: %w", err))
	}

	err = t.Execute(output, opts.Params.Data)

	if err != nil {
		panic(fmt.Errorf("Failed to execute template: %w", err))
	}

	outputPath, _ := outputLoc.Path()
	fmt.Println(templatePath, "->", outputPath)
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
