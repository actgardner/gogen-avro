package schema

import (
	"fmt"
	"github.com/actgardner/gogen-avro/generator"
)

const byteWriterInterface = `
type ByteWriter interface {
	Grow(int)
	WriteByte(byte) error
} 
`

const writeBoolMethod = `
func writeBool(r bool, w io.Writer) error {
	var b byte
	if r {
		b = byte(1)
	}

	var err error
	if bw, ok := w.(ByteWriter); ok {
		err = bw.WriteByte(b)
	} else {
		bb := make([]byte, 1)
		bb[0] = b
		_, err = w.Write(bb)
	}
	if err != nil {
		return err
	}
	return nil
}
`

type BoolField struct {
	PrimitiveField
}

func NewBoolField(definition interface{}) *BoolField {
	return &BoolField{PrimitiveField{
		definition:       definition,
		name:             "Bool",
		goType:           "bool",
		serializerMethod: "writeBool",
	}}
}

func (s *BoolField) AddSerializer(p *generator.Package) {
	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddFunction(UTIL_FILE, "", "writeBool", writeBoolMethod)
	p.AddImport(UTIL_FILE, "io")
}

func (s *BoolField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(bool); !ok {
		return "", fmt.Errorf("Expected bool as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}

func (s *BoolField) IsReadableBy(f AvroType) bool {
	_, ok := f.(*BoolField)
	return ok
}
