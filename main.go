package main

import (
	"fmt"
	"github.com/alanctgardner/gogen-avro/generator"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: gogen-avro <target directory> <schema files>\n")
		os.Exit(1)
	}
	targetDir := os.Args[1]
	files := os.Args[2:]
	pkg := generator.NewPackage("avro")
	for _, fileName := range files {
		schema, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %q - %v\n", fileName, err)
			os.Exit(2)
		}
		err = generator.DeserializeRecordSchema("avro", schema, pkg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating schema for file %q - %v\n", fileName, err)
			os.Exit(3)
		}
	}
	err := pkg.WriteFiles(targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing source files to directory %q - %v\n", targetDir, err)
		os.Exit(4)
	}
}
