package utils

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"regexp"

	"github.com/veigaribo/qveen/prompts"
)

var urlRegexp *regexp.Regexp

func init() {
	var err error
	urlRegexp, err = regexp.Compile(`^https?:\/\/`)

	if err != nil {
		panic(err)
	}
}

func IsStd(path string) bool {
	return path == "-"
}

func IsUrl(path string) bool {
	return urlRegexp.MatchString(path)
}

func IsLocal(path string) bool {
	return !IsStd(path) && !IsUrl(path)
}

func IsExplicitDir(path string) bool {
	end := path[len(path)-1]
	return (end == os.PathSeparator || end == '/') && IsLocal(path)
}

func OpenFileOrUrl(path string) (io.Reader, error) {
	if path == "" {
		return nil, errors.New("Tried to open an empty file path")
	}

	if IsStd(path) {
		return os.Stdin, nil // Since it's for reading, assume stdin
	}

	if IsUrl(path) {
		resp, err := http.Get(path)

		if err != nil {
			return nil, err
		}

		return resp.Body, nil
	}

	// Tautology.
	if IsLocal(path) {
		return os.Open(path)
	}

	return nil, nil // Unreachable.
}

var ErrAborted = errors.New("Aborted by user")

func FileWriter(path string, skipConfirm bool) (io.Writer, error) {
	if path == "" {
		return nil, errors.New("Tried to write to an empty file path")
	}

	if IsLocal(path) {
		var stat fs.FileInfo
		var err error

		if skipConfirm {
			goto ok
		}

		// Check for overwrites.
		stat, err = os.Stat(path)

		if err == nil {
			if stat.IsDir() {
				return nil, fmt.Errorf("Destination '%s' already exists and is a directory", path)
			}

			overwrite := prompts.AskConfirm(
				fmt.Sprintf("File '%s' already exists. Overwrite?", path),
			)

			if overwrite {
				goto ok
			} else {
				return nil, ErrAborted
			}
		}

	ok:
		return os.Create(path)
	} else if IsStd(path) {
		return os.Stdout, nil // Since it's for writing, assume stdout
	} else {
		return nil, fmt.Errorf("Tried to write to directory '%s' as if it were a file", path)
	}
}
