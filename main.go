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
		// TODO: Long description.

		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			templates.LeftDelim = leftDelimFlag
			templates.RightDelim = rightDelimFlag

			opts := RenderOptions{
				ParamsPath:   args[0],
				TemplatePath: templatePathFlag,
				OutputPath:   outputPathFlag,
				MetaKey:      metaKeyFlag,
				PromptValues: promptValueFlags,
				Overwrite:    overwriteFlag,
			}

			Render(opts)
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
