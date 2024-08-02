package params

import (
	"strings"
)

func GuessFormat(path string) *ParamsFormat {
	var format ParamsFormat

	dotI := strings.LastIndexByte(path, byte('.'))

	if dotI == -1 {
		return nil
	}

	// +1 to skip the dot itself.
	ext := path[dotI+1:]

	switch ext {
	case "toml":
		format = ParamsTomlFormat
		return &format
	case "json":
		fallthrough
	case "yaml":
		format = ParamsYamlFormat
		return &format
	}

	return nil
}
