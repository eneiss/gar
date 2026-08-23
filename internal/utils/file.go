package utils

import (
	"errors"
	"fmt"
	"os"
)

func FilesExist(files []string) error {
	for _, file := range files {
		if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("file does not exist: %s", file)
		}
	}
	return nil
}
