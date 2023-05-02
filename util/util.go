package util

import (
	"fmt"
)


func Processor(data []byte) {
	fmt.Printf("output -> \n %s", string(data))
}
 
func CommandLineGenerator(in, out string) []string {
    return []string{"rembg", "i", "-om", in,out} }
