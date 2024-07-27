package utils

func FirstOf[T comparable](options ...T) T {
	var empty T

	for _, option := range options {
		if option != empty {
			return option
		}
	}

	return empty
}
