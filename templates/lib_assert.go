package templates

func TemplateIsMap(val any) bool {
	_, ok := val.(map[any]any)
	return ok
}

func TemplateIsStr(val any) bool {
	_, ok := val.(string)
	return ok
}

func TemplateIsInt(val any) bool {
	_, ok := val.(int)
	return ok
}

// Trying to use a language agnostic idiom.
func TemplateIsArr(val any) bool {
	_, ok := val.([]any)
	return ok
}
