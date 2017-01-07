package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alanctgardner/gogen-avro/container"
	"github.com/alanctgardner/gogen-avro/generator"
	"github.com/alanctgardner/gogen-avro/types"
)

func main() {
	generateContainer := flag.Bool("container", false, "Whether to emit container file writer code")
	packageName := flag.String("package", "avro", "Name of generated package")
	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "Usage: gogen-avro [--container] [--package=<package name>] <target directory> <schema files>\n")
		os.Exit(1)
	}
	targetDir := flag.Arg(0)
	files := flag.Args()[1:]

	var err error
	pkg := generator.NewPackage(*packageName)

	if *generateContainer {
		err = addRecordDefinition([]byte(container.AVRO_BLOCK_SCHEMA), pkg, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating Avro container block schema - %v\n", err)
			os.Exit(2)
		}

		err = addRecordDefinition([]byte(container.AVRO_HEADER_SCHEMA), pkg, false)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating Avro container header schema - %v\n", err)
			os.Exit(2)
		}
	}

	for _, fileName := range files {
		schema, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %q - %v\n", fileName, err)
			os.Exit(2)
		}

		err = addRecordDefinition(schema, pkg, *generateContainer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error decoding schema for file %q - %v\n", fileName, err)
			os.Exit(3)
		}
	}

	err = pkg.WriteFiles(targetDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing source files to directory %q - %v\n", targetDir, err)
		os.Exit(4)
	}
}

func addRecordDefinition(schema []byte, pkg *generator.Package, generateContainer bool) error {
	recordDefinition, err := types.RecordDefinitionForSchema(schema)
	if err != nil {
		return err
	}
	recordDefinition.AddStruct(pkg)
	recordDefinition.AddSerializer(pkg)
	recordDefinition.AddDeserializer(pkg)

	if generateContainer {
		containerWriter := container.NewAvroContainerWriter(schema, recordDefinition)
		containerWriter.AddAvroContainerWriter(pkg)
	}
	return nil
}
