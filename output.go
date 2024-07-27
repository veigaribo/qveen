package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/veigaribo/qveen/prompts"
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

var ErrNoLeaf = errors.New("Output filename has only prefix.")
var ErrAborted = errors.New("Aborted by user.")

func (l *OutputLocation) Path() (string, error) {
	if l.Leaf == "" {
		return "", ErrNoLeaf
	}

	if strings.HasPrefix(l.Leaf, "/") {
		return l.Leaf, nil
	}

	return path.Join(l.Stem, l.Leaf), nil
}

func (l *OutputLocation) Writer() (io.Writer, error) {
	p, err := l.Path()

	if err != nil {
		return nil, err
	}

	stat, err := os.Stat(p)

	if err == nil {
		if stat.IsDir() {
			goto ok
		}

		overwrite := prompts.AskConfirm(fmt.Sprintf("File '%s' already exists. Overwrite?", p))

		if overwrite {
			goto ok
		} else {
			return nil, ErrAborted
		}
	}

ok:
	return os.Create(p)
}
