package main

import (
	"fmt"
	"github.com/alanctgardner/gogen-avro/example/avro"
	"os"
)

/* This example shows serializing records in a object container file */
func main() {
	// Create a new DemoSchema struct
	demoStruct := avro.DemoSchema{
		IntField:    1,
		DoubleField: 2.3,
		StringField: "A string",
		BoolField:   true,
		BytesField:  []byte{1, 2, 3, 4},
	}

	// Open a file to write
	fileWriter, err := os.Create("example_avro_container.avro")
	if err != nil {
		fmt.Printf("Error opening file writer: %v\n", err)
		return
	}

	// Create a DemoSchemaContainerWriter which writes to the file
	// Using the Null codec means blocks are uncompressed - other options are Snappy and Deflate
	var containerWriter *avro.DemoSchemaContainerWriter
	containerWriter, err = avro.NewDemoSchemaContainerWriter(fileWriter, avro.Null, 10)
	if err != nil {
		fmt.Printf("Error opening container writer: %v\n", err)
		return
	}

	// Write the record to the container file
	err = containerWriter.WriteRecord(demoStruct)
	if err != nil {
		fmt.Printf("Error writing record to file: %v\n", err)
		return
	}

	// Flush the buffers to ensure the last block has been written
	err = containerWriter.Flush()
	if err != nil {
		fmt.Printf("Error flushing last block to file: %v\n", err)
		return
	}

	fileWriter.Close()
}
