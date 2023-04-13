package util

import (
	"os"
)

func Cleanup(path string) error {
	return os.Remove(path)
}
