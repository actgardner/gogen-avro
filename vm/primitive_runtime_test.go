package vm

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/linkedin/goavro"
	"github.com/stretchr/testify/assert"
)

type PrimitiveStruct struct {
	One string
	Two int32
}

const PrimitiveSchema = `
{
  "type": "record",
  "name": "test",
  "fields": [
    {"name": "two", "type":"int"},
    {"name": "one", "type":"string"}
  ]
}
`

func (p *PrimitiveStruct) SetBoolean(field int, v bool) {}
func (p *PrimitiveStruct) SetInt(field int, v int32) {
	switch field {
	case 1:
		p.Two = v
		return
	}
	panic(fmt.Sprintf("Unexpected assignment to field %v with %v", field, v))
}
func (p *PrimitiveStruct) SetLong(field int, v int64)     {}
func (p *PrimitiveStruct) SetFloat(field int, v float32)  {}
func (p *PrimitiveStruct) SetDouble(field int, v float64) {}
func (p *PrimitiveStruct) SetBytes(field int, v []byte)   {}
func (p *PrimitiveStruct) SetString(field int, v string) {
	switch field {
	case 0:
		p.One = v
		return
	}
	panic(fmt.Sprintf("Unexpected assignment to field %v with %v", field, v))
}
func (p *PrimitiveStruct) Init(field int) {}
func (p *PrimitiveStruct) Get(field int) Assignable {
	return nil
}

func TestEvalPrimitive(t *testing.T) {
	program := []Instruction{
		Instruction{Op: Read, Type: Int, Field: 0},
		Instruction{Op: Set, Type: 3, Field: 1},
		Instruction{Op: Read, Type: String, Field: 0},
		Instruction{Op: Set, Type: 8, Field: 0},
	}

	codec, err := goavro.NewCodec(PrimitiveSchema)
	assert.Nil(t, err)

	encoded, err := codec.BinaryFromNative(nil, map[string]interface{}{
		"one": "hi",
		"two": 1234,
	})
	assert.Nil(t, err)
	reader := bytes.NewBuffer(encoded)
	var target PrimitiveStruct
	err = Eval(reader, program, &target)
	assert.Nil(t, err)
	assert.Equal(t, PrimitiveStruct{"hi", 1234}, target)
}

func BenchmarkPrimitiveDecode(b *testing.B) {
	codec, err := goavro.NewCodec(PrimitiveSchema)
	assert.Nil(b, err)

	encoded, err := codec.BinaryFromNative(nil, map[string]interface{}{
		"one": "hi",
		"two": 1234,
	})
	assert.Nil(b, err)

	program := []Instruction{
		Instruction{Op: Read, Type: Int, Field: 0},
		Instruction{Op: Set, Type: 3, Field: 1},
		Instruction{Op: Read, Type: String, Field: 0},
		Instruction{Op: Set, Type: 8, Field: 0},
	}

	for i := 0; i < b.N; i++ {
		reader := bytes.NewBuffer(encoded)
		var target PrimitiveStruct
		err = Eval(reader, program, &target)
	}
}

func BenchmarkPrimitiveGoavro(b *testing.B) {
	codec, err := goavro.NewCodec(PrimitiveSchema)
	assert.Nil(b, err)

	encoded, err := codec.BinaryFromNative(nil, map[string]interface{}{
		"one": "hi",
		"two": 1234,
	})
	assert.Nil(b, err)

	for i := 0; i < b.N; i++ {
		codec.NativeFromBinary(encoded)
	}
}
