package prompts

import (
	"reflect"
	"strings"

	"github.com/charmbracelet/huh"
)

var SupportedPromptKinds = []string{
	"input", "text", "select", "confirm",
}

type Prompt struct {
	Kind  string
	Name  string
	Title string

	Specific any
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
	if len(prompts) == 0 {
		return make(map[string]any), nil
	}

	valuePtrs := make(map[string]any)
	var fields []huh.Field

	for _, prompt := range prompts {
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

	group := huh.NewGroup(fields...)
	form := huh.NewForm(group)

	values := make(map[string]any)
	err := form.Run()

	if err != nil {
		return values, err
	}

	for key, valuePtr := range valuePtrs {
		values[key] = reflect.ValueOf(valuePtr).Elem().Interface()
	}

	return values, nil
}

func promptInput(prompt Prompt, ptrs map[string]any) huh.Field {
	var value string
	title := prompt.Title

	if title == "" {
		title = defaultTitle(prompt.Name)
	}

	ptrs[prompt.Name] = &value
	return huh.NewInput().Title(title).Value(&value)
}

func promptText(prompt Prompt, ptrs map[string]any) huh.Field {
	var value string
	title := prompt.Title

	if title == "" {
		title = defaultTitle(prompt.Name)
	}

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

	title := prompt.Title
	if title == "" {
		title = defaultTitle(prompt.Name)
	}

	ptrs[prompt.Name] = &value

	return huh.NewSelect[any]().
		Title(title).
		Options(huhOptions...).
		Value(&value)
}

func promptConfirm(prompt Prompt, ptrs map[string]any) huh.Field {
	var value bool
	title := prompt.Title

	if title == "" {
		title = defaultTitle(prompt.Name)
	}

	ptrs[prompt.Name] = &value
	return huh.NewConfirm().Title(title).Value(&value)
}

func defaultTitle(name string) string {
	var builder strings.Builder
	builder.WriteString("Value for '")
	builder.WriteString(name)
	builder.WriteString("':")

	return builder.String()
}
