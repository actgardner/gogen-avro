package main

import (
	"fmt"
	"github.com/alanctgardner/avro-generator/generator"
	"io/ioutil"
	"os"
)

func main() {
	file := os.Args[1]
	schema, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file %q - %v", file, err)
		return
	}
	goDefs, err := generator.GenerateForSchema(schema)
	if err != nil {
		fmt.Printf("Error generating schema for file %q - %v", file, err)
		return
	}
	fmt.Printf(goDefs)
}
