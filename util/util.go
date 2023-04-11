package util

import (
	"io/fs"
	"os"
)

func Cleanup(path string) *os.PathError {
	return os.Remove(path).(*fs.PathError)
}
