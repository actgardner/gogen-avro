package vm

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/linkedin/goavro"
	"github.com/stretchr/testify/assert"
)

type PrimitiveArrayStruct struct {
	One   int32
	Two   []string
	Three int32
}

type StringArray []string

func (p *StringArray) SetBoolean(field int, v bool)   {}
func (p *StringArray) SetInt(field int, v int32)      {}
func (p *StringArray) SetLong(field int, v int64)     {}
func (p *StringArray) SetFloat(field int, v float32)  {}
func (p *StringArray) SetDouble(field int, v float64) {}
func (p *StringArray) SetBytes(field int, v []byte)   {}
func (p *StringArray) SetString(field int, v string) {
	*p = append(*p, v)
}
func (p *StringArray) Init(field int) {}
func (p *StringArray) Get(field int) Assignable {
	return nil
}

const PrimitiveArraySchema = `
{
  "type": "record",
  "name": "test",
  "fields": [
    {"name": "one", "type": "int"},
    {
      "name": "two", 
      "type": {
        "type": "array",
        "items": "string" 
      }
    },
    {"name": "three", "type": "int"}
  ]
}
`

func (p *PrimitiveArrayStruct) SetBoolean(field int, v bool) {}
func (p *PrimitiveArrayStruct) SetInt(field int, v int32) {
	switch field {
	case 0:
		p.One = v
		return
	case 2:
		p.Three = v
		return

	}
	panic(fmt.Sprintf("Unexpected assignment to field %v with %v", field, v))
}
func (p *PrimitiveArrayStruct) SetLong(field int, v int64)     {}
func (p *PrimitiveArrayStruct) SetFloat(field int, v float32)  {}
func (p *PrimitiveArrayStruct) SetDouble(field int, v float64) {}
func (p *PrimitiveArrayStruct) SetBytes(field int, v []byte)   {}
func (p *PrimitiveArrayStruct) SetString(field int, v string) {
	switch field {
	case 1:
		if p.Two == nil {
			p.Two = make(StringArray, 0)
		}
		(*StringArray)(&p.Two).SetString(0, v)
		return
	}
	panic(fmt.Sprintf("Unexpected assignment to field %v with %v", field, v))
}
func (p *PrimitiveArrayStruct) Init(field int) {}
func (p *PrimitiveArrayStruct) Get(field int) Assignable {
	return nil
}

func TestEvalPrimitiveArray(t *testing.T) {
	program := []Instruction{
		Instruction{Op: Read, Type: Int, Field: 65535},
		Instruction{Op: Set, Type: Int, Field: 0},
		Instruction{Op: BlockStart, Type: Unused, Field: 65535},
		Instruction{Op: Read, Type: String, Field: 65535},
		Instruction{Op: Set, Type: String, Field: 1},
		Instruction{Op: BlockEnd, Type: Unused, Field: 65535},
		Instruction{Op: Read, Type: Int, Field: 65535},
		Instruction{Op: Set, Type: Int, Field: 2},
	}

	codec, err := goavro.NewCodec(PrimitiveArraySchema)
	assert.Nil(t, err)

	encoded, err := codec.BinaryFromNative(nil, map[string]interface{}{
		"one":   789,
		"two":   []string{"a", "b", "c"},
		"three": 1234,
	})
	assert.Nil(t, err)
	reader := bytes.NewBuffer(encoded)
	var target PrimitiveArrayStruct
	err = Eval(reader, program, &target)
	assert.Nil(t, err)
	assert.Equal(t, PrimitiveArrayStruct{789, []string{"a", "b", "c"}, 1234}, target)
}
