package main

import (
	"fmt"
	"github.com/alanctgardner/gogen-avro/generator"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: gogen-avro <schema file> <target directory>\n")
		os.Exit(1)
	}
	file := os.Args[1]
	targetDir := os.Args[2]
	schema, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %q - %v\n", file, err)
		os.Exit(2)
	}
	pkg, err := generator.GenerateForSchema("avro", schema)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating schema for file %q - %v\n", file, err)
		os.Exit(3)
	}
	err = pkg.WriteFiles(targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing source files to directory %q - %v\n", targetDir, err)
		os.Exit(4)
	}
}
