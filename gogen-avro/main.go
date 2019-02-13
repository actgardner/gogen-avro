package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/actgardner/gogen-avro/generator"
	"github.com/actgardner/gogen-avro/types"
)

func main() {
	packageName := flag.String("package", "avro", "Name of generated package")
	containers := flag.Bool("containers", false, "Whether to generate container writer methods")
	shortUnions := flag.Bool("short-unions", false, "Whether to use shorter names for Union types")
	namespacedNames := flag.Bool("namespaced-names", false, "Whether to generate namespaced names for types")

	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "Usage: gogen-avro [--namespaced-names] [--short-unions] [--package=<package name>] [--containers] <target directory> <schema files>\n")
		os.Exit(1)
	}

	targetDir := flag.Arg(0)
	files := flag.Args()[1:]

	var err error
	pkg := generator.NewPackage(*packageName)
	namespace := types.NewNamespace(*shortUnions)

	if *namespacedNames {
		generator.SetNamer(generator.NewNamespaceNamer())
	}

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
	const fileComment = `// Code generated by github.com/actgardner/gogen-avro. DO NOT EDIT.
/*
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
