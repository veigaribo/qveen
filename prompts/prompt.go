package prompts

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"syscall"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/x/term"
)

var SupportedPromptKinds = []string{
	"input", "text", "select", "confirm",
}

type Prompt struct {
	Kind     string
	Name     string
	Title    string
	Specific any

	Value any
}

func (p *Prompt) TryPrefill(raw string) error {
	switch p.Kind {
	case "input":
		fallthrough
	case "text":
		p.Value = raw
	case "select":
		specific := p.Specific.(PromptSelectSpecific)

		for _, option := range specific.Options {
			if option.Title == raw {
				p.Value = option.Value
				return nil
			}
		}

		return errors.New("Given value for `select` prompt does not correspond to one of the options. (Needs to match the title)")
	case "confirm":
		value, err := strconv.ParseBool(raw)

		if err != nil {
			return err
		}

		p.Value = value
		return nil
	}

	return nil
}

func (p Prompt) GetTitle() string {
	if p.Title != "" {
		return p.Title
	}

	var builder strings.Builder
	builder.WriteString("Value for '")
	builder.WriteString(p.Name)
	builder.WriteString("':")

	return builder.String()
}

type PromptSelectOption struct {
	Title string
	Value any
}

type PromptSelectSpecific struct {
	Options []PromptSelectOption
}

func AskConfirm(title string) bool {
	var confirm bool

	huh.NewConfirm().
		Title(title).
		Affirmative("Yes").
		Negative("No").
		Value(&confirm).
		Run()

	return confirm
}

func DoPrompt(prompts []Prompt) (map[string]any, error) {
	var err error

	if len(prompts) == 0 {
		return make(map[string]any), nil
	}

	// TODO: Prompt values as flags.
	if !term.IsTerminal(uintptr(syscall.Stdin)) {
		return nil, errors.New("Tried to prompt while not connected to a terminal")
	}

	valuePtrs := make(map[string]any)
	values := make(map[string]any)
	var fields []huh.Field

	for _, prompt := range prompts {
		if prompt.Value != nil {
			// Prefilled.
			valuePtrs[prompt.Name] = &prompt.Value
			continue
		}

		switch prompt.Kind {
		case "input":
			fields = append(fields, promptInput(prompt, valuePtrs))
		case "text":
			fields = append(fields, promptText(prompt, valuePtrs))
		case "select":
			fields = append(fields, promptSelect(prompt, valuePtrs))
		case "confirm":
			fields = append(fields, promptConfirm(prompt, valuePtrs))
		}
	}

	var group *huh.Group
	var form *huh.Form

	if len(fields) == 0 {
		// If everything was prefilled.
		goto resolve
	}

	group = huh.NewGroup(fields...)
	form = huh.NewForm(group)

	err = form.Run()

	if err != nil {
		return values, err
	}

resolve:
	for key, valuePtr := range valuePtrs {
		values[key] = reflect.ValueOf(valuePtr).Elem().Interface()
	}

	return values, nil
}

func promptInput(prompt Prompt, ptrs map[string]any) huh.Field {
	var value string
	title := prompt.GetTitle()

	ptrs[prompt.Name] = &value
	return huh.NewInput().Title(title).Value(&value)
}

func promptText(prompt Prompt, ptrs map[string]any) huh.Field {
	var value string
	title := prompt.GetTitle()

	ptrs[prompt.Name] = &value
	return huh.NewText().Title(title).Value(&value)
}

func promptSelect(prompt Prompt, ptrs map[string]any) huh.Field {
	var value any

	specific := prompt.Specific.(PromptSelectSpecific)
	var huhOptions []huh.Option[any]

	for _, option := range specific.Options {
		huhOptions = append(huhOptions,
			huh.NewOption(option.Title, option.Value))
	}

	title := prompt.GetTitle()
	ptrs[prompt.Name] = &value

	return huh.NewSelect[any]().
		Title(title).
		Options(huhOptions...).
		Value(&value)
}

func promptConfirm(prompt Prompt, ptrs map[string]any) huh.Field {
	var value bool
	title := prompt.GetTitle()

	ptrs[prompt.Name] = &value
	return huh.NewConfirm().Title(title).Value(&value)
}
