// This example shows serializing and deserializing records in single-object encoding
package main

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/actgardner/gogen-avro/v10/example/avro"
	"github.com/actgardner/gogen-avro/v10/soe"
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

	buf := &bytes.Buffer{}

	// Write the record with single-object encoding to the buffer
	err := soe.WriteRecord(buf, &demoStruct)
	if err != nil {
		fmt.Printf("Error writing record to buffer: %v\n", err)
		return
	}

	// Read the single-object header from the buffer to get the schema fingerprint
	fingerprint, err := soe.ReadHeader(buf)
	if err != nil {
		fmt.Printf("Error reading single-object header: %v\n", err)
		return
	}

	// Check that the schema fingeprint matches the fingerprint we expect
	if string(fingerprint) != avro.DemoSchemaAvroCRC64Fingerprint {
		fmt.Printf("Object had unexpected fingerprint: %v\n", fingerprint)
		return
	}

	output, err := avro.DeserializeDemoSchema(buf)
	if err != nil {
		fmt.Printf("Error reading single-object header: %v\n", err)
		return
	}

	if !reflect.DeepEqual(output, demoStruct) {
		fmt.Printf("Expected structs to match! %v %v\n", output, demoStruct)
		return
	}
	fmt.Printf("Round-trip successful!")
}
