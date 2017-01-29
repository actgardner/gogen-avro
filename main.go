package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/alanctgardner/gogen-avro/generator"
	"github.com/alanctgardner/gogen-avro/types"
)

func main() {
	packageName := flag.String("package", "avro", "Name of generated package")
	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "Usage: gogen-avro [--package=<package name>] <target directory> <schema files>\n")
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

		_, err = namespace.FieldDefinitionForSchema(schema)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding schema for file %q - %v\n", fileName, err)
			os.Exit(3)
		}
	}

	// Resolve dependencies and add the schemas to the package
	err = addFieldsToPackage(namespace, pkg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating code for schema - %v\n", err)
		os.Exit(4)
	}

	// Add header comment to all generated files.
	for _, f := range pkg.Files() {
		pkg.AddHeader(f, codegenComment(files))
	}

	err = pkg.WriteFiles(targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing source files to directory %q - %v\n", targetDir, err)
		os.Exit(4)
	}
}

func addFieldsToPackage(namespace *types.Namespace, pkg *generator.Package) error {
	for _, schema := range namespace.Schemas {
		err := schema.Root.ResolveReferences(namespace)
		if err != nil {
			return err
		}

		schema.Root.AddStruct(pkg)
		schema.Root.AddSerializer(pkg)
		schema.Root.AddDeserializer(pkg)
	}
	return nil
}

// codegenComment generates a comment informing readers they are looking at
// generated code and lists the source avro files used to generate the code
//
// invariant: sources > 0
func codegenComment(sources []string) string {
	const fileComment = `/*
 * CODE GENERATED AUTOMATICALLY WITH github.com/alanctgardner/gogen-avro
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
