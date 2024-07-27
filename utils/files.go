package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
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

func IsLocalFile(path string) bool {
	return !IsStd(path) && !IsUrl(path)
}

func IsExplicitDir(path string) bool {
	return strings.HasSuffix(path, "/") && IsLocalFile(path)
}

func OpenFileOrUrl(path string) (io.Reader, error) {
	if path == "" {
		return nil, errors.New("Tried to open an empty file path")
	}

	if IsUrl(path) {
		resp, err := http.Get(path)

		if err != nil {
			return nil, err
		}

		return resp.Body, nil
	} else {
		if IsLocalFile(path) {
			return os.Open(path)
		} else if IsStd(path) {
			return os.Stdin, nil // Since it's for reading, assume stdin
		} else {
			return nil, fmt.Errorf("Tried to open directory '%s' like a file", path)
		}
	}
}

func FileWriter(path string) (io.Writer, error) {
	if path == "" {
		return nil, errors.New("Tried to write to an empty file path")
	}

	if IsLocalFile(path) {
		return os.Create(path)
	} else if IsStd(path) {
		return os.Stdout, nil // Since it's for writing, assume stdout
	} else {
		return nil, fmt.Errorf("Tried to write to directory '%s' as if it were a file", path)
	}
}
