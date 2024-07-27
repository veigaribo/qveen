package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/veigaribo/qveen/templates"
)

func catchPanic() {
	e := recover()

	if e != nil {
		fmt.Fprintln(os.Stderr, e)
	}
}

type FlagType uint

const (
	StringFlagType FlagType = iota
	BoolFlagType
	StringToStringType
)

func (typ FlagType) AllowsMultiple() bool {
	return typ == StringToStringType
}

type Flag struct {
	Type          FlagType
	Short         string
	Long          string
	ParameterName string
	Target        any // Really a pointer.
	Description   string
}

func main() {
	debug := os.Getenv("DEBUG")

	if debug == "" {
		// Avoid printing stack trace by default.
		defer catchPanic()
	}

	var templatePathFlag string
	var outputPathFlag string
	var promptValueFlags map[string]string

	var metaKeyFlag string
	var leftDelimFlag string
	var rightDelimFlag string
	var overwriteFlag bool

	rootCmd := cobra.Command{
		Use:   "qveen",
		Short: "Generate files from templates.",
		Long: `Generate files from templates.

The parameter file should contain valid TOML. Its
contents will be used as the data from which to render the template.
If a value of ` + "`-`" + ` is given, the contents will be expected to come
from stdin.

A [meta] table in the parameter file is treated in a special way.
It may contain the following keys:

- template: Go template file
- output: File to output after rendering
- prompt: Values to be provided interactively

The prompt key is expected to contain an array of tables with the
following keys:

- kind: Determines the type of prompt to present
- name: Name of the variable in which to bind
- title: Text to show when prompting

Currently, the allowed values for ` + "`kind`" + ` are { text }.

If a variable is present both in the ` + "`meta.prompt`" + ` array and as a
standalone value, the former one will be preferred.

All of the above values may reference other values defined in the
file using regular template syntax. However, expansion is not recursive:
If a field containing a placeholder is included into another, that
placeholder will be kept in the final value as is.

The ` + "`template`" + ` and ` + "`output`" + ` flags will be used in case
they are missing from the parameter file if and only if a single
parameter file is provided. In the case of ` + "`output`" + `, if one of the
values ends with /, it its considered a prefix to apply to the other one.`,
		Args: cobra.MinimumNArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			templates.LeftDelim = leftDelimFlag
			templates.RightDelim = rightDelimFlag

			if len(args) == 1 {
				opts := Render1Options{
					ParamsPath:   args[0],
					TemplatePath: templatePathFlag,
					OutputPath:   outputPathFlag,
					MetaKey:      metaKeyFlag,
					PromptValues: promptValueFlags,
					Overwrite:    overwriteFlag,
				}

				Render1(opts)
			} else {
				if templatePathFlag != "" {
					fmt.Fprintln(os.Stderr, "Ignoring `template` flag on multi-file mode.")
				}

				opts := RenderNOptions{
					ParamsPaths:   args,
					OutputDirPath: outputPathFlag,
					MetaKey:       metaKeyFlag,
					PromptValues:  promptValueFlags,
					Overwrite:     overwriteFlag,
				}

				RenderN(opts)
			}
		},
	}

	// We use this array to build with "usage" section when helping, and
	// also as a base to actually register the flags with `cobra`.
	flags := []Flag{
		{
			Type:          StringFlagType,
			Short:         "t",
			Long:          "template",
			ParameterName: "template-file",
			Target:        &templatePathFlag,
			Description:   "Go template file to use (placeholders allowed).",
		},
		{
			Type:          StringFlagType,
			Short:         "o",
			Long:          "output",
			ParameterName: "output-file",
			Target:        &outputPathFlag,
			Description:   "Destination file name (placeholders allowed).",
		},
		{
			Type:          StringToStringType,
			Short:         "p",
			Long:          "prompt-value",
			ParameterName: "key=val",
			Target:        &promptValueFlags,
			Description:   "Sets a value for a prompt upfront.",
		},
		{
			Type:          StringFlagType,
			Short:         "m",
			Long:          "meta",
			ParameterName: "meta-key",
			Target:        &metaKeyFlag,
			Description:   "Key to look for meta information, instead of `meta`.",
		},
		{
			Type:          StringFlagType,
			Short:         "l",
			Long:          "left-delim",
			ParameterName: "delim",
			Target:        &leftDelimFlag,
			Description:   "String to use as the left delimiter for the template instead of `{{`.",
		},
		{
			Type:          StringFlagType,
			Short:         "r",
			Long:          "right-delim",
			ParameterName: "delim",
			Target:        &rightDelimFlag,
			Description:   "String to use as the right delimiter for the template instead of `}}`.",
		},
		{
			Type:        BoolFlagType,
			Short:       "y",
			Long:        "overwrite",
			Target:      &overwriteFlag,
			Description: "If set, won't ask for confirmation when overwriting files.",
		},
	}

	for _, flag := range flags {
		registerFlag(&rootCmd, flag)
	}

	// Add information for `--help` so it shows up in the usage, but
	// don't really add it since it is created automatically.
	var null bool
	flags = append(flags, Flag{
		Type:        BoolFlagType,
		Short:       "h",
		Long:        "help",
		Target:      &null,
		Description: "Show this message and exit.",
	})

	rootCmd.SetUsageFunc(usage(flags))

	err := rootCmd.Execute()

	if err != nil {
		os.Exit(1)
	}
}

func registerFlag(cmd *cobra.Command, flag Flag) {
	switch flag.Type {
	case StringFlagType:
		target := flag.Target.(*string)

		cmd.Flags().StringVarP(
			target,
			flag.Long,
			flag.Short,
			"",
			flag.Description,
		)
	case BoolFlagType:
		target := flag.Target.(*bool)

		cmd.Flags().BoolVarP(
			target,
			flag.Long,
			flag.Short,
			false,
			flag.Description,
		)
	case StringToStringType:
		target := flag.Target.(*map[string]string)

		cmd.Flags().StringToStringVarP(
			target,
			flag.Long,
			flag.Short,
			make(map[string]string),
			flag.Description,
		)
	}
}

func usage(flags []Flag) func(cmd *cobra.Command) error {
	return func(cmd *cobra.Command) error {
		// Not a hard limit. Will not break line if not in a good
		// position to do so. A bit lower than usual to compensate.
		const maxLineLength = 70

		var builder strings.Builder

		// Characters since last line break.
		// Used for wrapping.
		rowLen := 0

		// Assumes no line break inside `str`!
		writeLine := func(str string) {
			builder.WriteString(str)
			rowLen += len(str)
		}

		breakLine := func() {
			builder.WriteRune('\n')
			rowLen = 0
		}

		writeLine("Usage:")
		breakLine()
		writeLine("  ")

		writeLine(os.Args[0])
		writeLine(" ")

		// Column at which flags start.
		// Basically makes indentation go here on line break when
		// using `maybeLineBreak` specifically.
		baseColumn := rowLen - 1

		writeFlag := func(flag Flag) {
			writeLine("[-")
			writeLine(flag.Short)
			writeLine(" | --")
			writeLine(flag.Long)

			if flag.ParameterName != "" {
				writeLine(" <")
				writeLine(flag.ParameterName)
				writeLine(">")
			}

			if flag.Type.AllowsMultiple() {
				builder.WriteString(" ...")
			}

			writeLine("]")
		}

		maybeBreakLine := func() bool {
			if rowLen >= maxLineLength-1 {
				breakLine()

				for rowLen < baseColumn {
					writeLine(" ")
				}

				return true
			}

			return false
		}

		flag := flags[0]
		writeFlag(flag)
		maybeBreakLine()

		for _, flag := range flags[1:] {
			writeLine(" ")
			writeFlag(flag)
			maybeBreakLine()
		}

		writeLine(" <params-file>")
		maybeBreakLine()
		writeLine(" [<params-file>...]")
		breakLine()
		breakLine()

		writeLine("Options:")
		breakLine()

		optionHeader := func(flag Flag) string {
			var builder strings.Builder
			builder.WriteString("-")
			builder.WriteString(flag.Short)
			builder.WriteString(", --")
			builder.WriteString(flag.Long)

			return builder.String()
		}

		descriptionHeaders := make([]string, 0, len(flags))

		// Column at which description starts.
		baseColumn = 0

		for _, flag := range flags {
			header := optionHeader(flag)
			descriptionHeaders = append(descriptionHeaders, header)

			if len(header) > baseColumn {
				baseColumn = len(header)
			}
		}

		baseColumn += 3 // Account for spaces we put at the start.

		for i, flag := range flags {
			writeLine("  ")
			writeLine(descriptionHeaders[i])

			for rowLen < baseColumn {
				writeLine(" ")
			}

			words := strings.Fields(flag.Description)
			writeLine(words[0])

			for _, word := range words[1:] {
				broke := maybeBreakLine()

				if !broke {
					writeLine(" ")
				}

				writeLine(word)
			}

			breakLine()
		}

		breakLine()
		fmt.Fprint(cmd.OutOrStderr(), builder.String())
		return nil
	}
}
