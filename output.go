package main

import (
	"errors"
	"path"
	"strings"

	"github.com/veigaribo/qveen/utils"
)

// Constructs a path with optional prefix.
type OutputLocation struct {
	Stem string
	Leaf string
}

func IsPrefix(filename string) bool {
	return strings.HasSuffix(filename, "/")
}

func (l *OutputLocation) Add(filename string) {
	if len(filename) == 0 {
		return
	}

	if IsPrefix(filename) {
		l.Stem = filename
	} else {
		l.Leaf = filename
	}
}

var ErrNoLeaf = errors.New("Output filename has only prefix")

func (l *OutputLocation) Path() (string, error) {
	if l.Leaf == "" {
		return "", ErrNoLeaf
	}

	if utils.IsStd(l.Leaf) {
		return l.Leaf, nil
	}

	return path.Join(l.Stem, l.Leaf), nil
}
