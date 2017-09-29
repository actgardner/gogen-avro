package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/alanctgardner/gogen-avro.v5/generator"
	"gopkg.in/alanctgardner/gogen-avro.v5/types"
)

func main() {
	packageName := flag.String("package", "avro", "Name of generated package")
	containers := flag.Bool("containers", false, "Whether to generate container writer methods")

	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "Usage: gogen-avro [--package=<package name>] [--containers] <target directory> <schema files>\n")
		os.Exit(1)
	}

	targetDir := flag.Arg(0)
	files := flag.Args()[1:]

	var err error
	pkg := generator.NewPackage(*packageName)
	namespace := types.NewNamespace()

	for _, fileName := range files {
		schema, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %q - %v\n", fileName, err)
			os.Exit(2)
		}

		_, err = namespace.TypeForSchema(schema)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding schema for file %q - %v\n", fileName, err)
			os.Exit(3)
		}
	}

	err = namespace.AddToPackage(pkg, codegenComment(files), *containers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating code for schema - %v\n", err)
		os.Exit(4)
	}

	err = pkg.WriteFiles(targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing source files to directory %q - %v\n", targetDir, err)
		os.Exit(4)
	}
}

// codegenComment generates a comment informing readers they are looking at
// generated code and lists the source avro files used to generate the code
//
// invariant: sources > 0
func codegenComment(sources []string) string {
	const fileComment = `/*
 * CODE GENERATED AUTOMATICALLY WITH gopkg.in/alanctgardner/gogen-avro.v5
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 *
 * %s
 */`
	var sourceBlock []string
	if len(sources) == 1 {
		sourceBlock = append(sourceBlock, "SOURCE:")
	} else {
		sourceBlock = append(sourceBlock, "SOURCES:")
	}

	for _, source := range sources {
		_, fName := filepath.Split(source)
		sourceBlock = append(sourceBlock, fmt.Sprintf(" *     %s", fName))
	}

	return fmt.Sprintf(fileComment, strings.Join(sourceBlock, "\n"))
}
