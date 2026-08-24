package utils

import (
	"errors"
	"fmt"
	"os"
)

// FilesExist checks if all paths in the files input exist on the current
// filesystem.
// It returns an error on the first non-existent file found, or nil if all
// files exist.
func FilesExist(files []string) error {
	for _, file := range files {
		return FileExists(file)
	}
	return nil
}

// FileExists checks if the provided argument corresponds to an existing file
// on the current filesystem. Returns nil if the file exists, or an error if it
// doesn't.
func FileExists(file string) error {
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file does not exist: %s", file)
	}
	return nil
}
