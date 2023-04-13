package util

import (
	"fmt"
	"os"
)

func Cleanup(path string) error {
	return os.Remove(path)
}

func Processor(data []byte) {
	fmt.Printf("output -> \n %s", string(data))
}
