package main

import (
	"fmt"
	"github.com/alanctgardner/gogen-avro/generator"
	"io/ioutil"
	"os"
)

func main() {
	file := os.Args[1]
	schema, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %q - %v\n", file, err)
		os.Exit(1)
	}
	goDefs, err := generator.GenerateForSchema(schema)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating schema for file %q - %v\n", file, err)
		os.Exit(1)
	}
	fmt.Printf(goDefs)
}
