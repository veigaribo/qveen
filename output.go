package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/veigaribo/qveen/prompts"
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
var ErrAborted = errors.New("Aborted by user")

func (l *OutputLocation) Path() (string, error) {
	if l.Leaf == "" {
		return "", ErrNoLeaf
	}

	if utils.IsStd(l.Leaf) {
		return l.Leaf, nil
	}

	p := path.Join(l.Stem, l.Leaf)
	stat, err := os.Stat(p)

	if err == nil {
		if stat.IsDir() {
			goto ok
		}

		overwrite := prompts.AskConfirm(
			fmt.Sprintf("File '%s' already exists. Overwrite?", p),
		)

		if overwrite {
			goto ok
		} else {
			return "", ErrAborted
		}
	}

ok:
	return p, nil
}
