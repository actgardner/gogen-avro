package main

import (
	"bytes"
	"fmt"
	"github.com/actgardner/gogen-avro/example/avro"
)

// This example shows serializing and deserializing records as byte buffers
func main() {
	// Create a new DemoSchema struct
	demoStruct := avro.DemoSchema{
		IntField:    1,
		DoubleField: 2.3,
		StringField: "A string",
		BoolField:   true,
		BytesField:  []byte{1, 2, 3, 4},
	}

	// Serialize the struct to a byte buffer
	var buf bytes.Buffer
	fmt.Printf("Serializing struct: %#v\n", demoStruct)
	demoStruct.Serialize(&buf)

	// Deserialize the byte buffer back into a struct
	newDemoStruct, err := avro.DeserializeDemoSchema(&buf)
	if err != nil {
		fmt.Printf("Error deserializing struct: %v\n", err)
		return
	}
	fmt.Printf("Deserialized struct: %#v\n", newDemoStruct)
}
