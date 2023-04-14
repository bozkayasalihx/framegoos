package util

import (
	"fmt"
	"os"
)

func Cleanup(path string) error {
	return os.RemoveAll(path)
}

func Processor(data []byte) {
	fmt.Printf("output -> \n %s", string(data))
}
