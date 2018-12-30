package vm

import (
	"bytes"
	"testing"

	"github.com/actgardner/gogen-avro/types"

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

func (p *PrimitiveStruct) SetBoolean(v bool)   {}
func (p *PrimitiveStruct) SetInt(v int32)      {}
func (p *PrimitiveStruct) SetLong(v int64)     {}
func (p *PrimitiveStruct) SetFloat(v float32)  {}
func (p *PrimitiveStruct) SetDouble(v float64) {}
func (p *PrimitiveStruct) SetBytes(v []byte)   {}
func (p *PrimitiveStruct) SetString(v string)  {}
func (p *PrimitiveStruct) Get(field int) types.Field {
	switch field {
	case 0:
		return (*types.String)(&p.One)
	case 1:
		return (*types.Int)(&p.Two)
	}
	panic("Field index out of range!")
}
func (p *PrimitiveStruct) AppendMap(v string) types.Field { return nil }
func (p *PrimitiveStruct) AppendArray() types.Field       { return nil }

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
