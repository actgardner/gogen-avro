// This example shows serializing and deserializing records in a object container file
package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/actgardner/gogen-avro/v7/container"
	"github.com/actgardner/gogen-avro/v7/example/avro"
)

func main() {
	// Create a new DemoSchema struct
	demoStruct := avro.DemoSchema{
		IntField:    1,
		DoubleField: 2.3,
		StringField: "A string",
		BoolField:   true,
		BytesField:  []byte{1, 2, 3, 4},
	}

	// Write to an in-memory buffer. You could use os.Create
	// to write to an actual file.
	buffer := new(bytes.Buffer)

	// Create a container.Writer which can write any generated Avro struct to a file
	// Note that all the objects written to the file must be the same type
	// Using the Null codec means blocks are uncompressed - other options are Snappy and Deflate
	containerWriter, err := container.NewWriter(buffer, container.Null, 10, demoStruct.Schema())
	if err != nil {
		fmt.Printf("Error opening container writer: %v\n", err)
		return
	}

	// Write the record to the container file
	if err := containerWriter.WriteRecord(&demoStruct); err != nil {
		fmt.Printf("Error writing record to file: %v\n", err)
		return
	}

	// Flush the buffers to ensure the last block has been written
	if err := containerWriter.Flush(); err != nil {
		fmt.Printf("Error flushing last block to file: %v\n", err)
		return
	}

	// Read the data again. You could use os.Open to read from a file.
	fileReader := bytes.NewReader(buffer.Bytes())

	// Create a new OCF reader
	ocfReader, err := container.NewReader(fileReader)
	if err != nil {
		fmt.Printf("Error creating OCF file reader: %v\n", err)
		return
	}

	// Read the records back until the file is finished
	for {
		record, err := avro.DeserializeDemoSchema(ocfReader)
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Printf("Error deserializing record: %v\n", err)
			return
		}
		fmt.Printf("Read record: %#v\n", record)
	}
}
