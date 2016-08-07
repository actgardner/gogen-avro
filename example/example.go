package main

import (
	"bytes"
	"fmt"
	"github.com/alanctgardner/gogen-avro/example/avro"
)

// Use a go:generate directive to build the Go structs for `example.avsc`
// Source files will be in a package called `avro`

//go:generate $GOPATH/bin/gogen-avro ./avro example.avsc

func main() {
	// Create a new DemoSchema struct
	demoStruct := &avro.DemoSchema{
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
