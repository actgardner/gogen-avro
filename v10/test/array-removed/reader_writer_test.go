package avro

import (
	"bytes"
	"github.com/actgardner/gogen-avro/v10/compiler"
	"github.com/actgardner/gogen-avro/v10/test/array-removed/reader"
	"github.com/actgardner/gogen-avro/v10/test/array-removed/writer"
	"github.com/actgardner/gogen-avro/v10/vm"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestReaderWriter(t *testing.T) {

	writerSchema, err := os.ReadFile("writer.avsc")
	if err != nil {
		panic(err)
	}

	readerSchema, err := os.ReadFile("reader.avsc")
	if err != nil {
		panic(err)
	}

	prog, err := compiler.CompileSchemaBytes(writerSchema, readerSchema)
	if err != nil {
		panic(err)
	}

	sourceRecord := writer.Foo{
		F0: []writer.Bar{
			{16},
			{32},
			{64},
		},
		F1: []int32{1, 2, 3},
		F2: []string{"a", "b"},
	}

	output := bytes.NewBuffer(make([]byte, 0))
	if err := sourceRecord.Serialize(output); err != nil {
		panic(err)
	}

	input := bytes.NewReader(output.Bytes())

	readRecord := reader.NewFoo()

	if err := vm.Eval(input, prog, &readRecord); err != nil {
		panic(err)
	}

	assert.Equal(t,
		reader.Foo{
			F2: []string{"a", "b"},
		},
		readRecord)
}
