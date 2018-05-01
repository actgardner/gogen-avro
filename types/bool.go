package types

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

const byteReaderInterface = `
type ByteReader interface {
	ReadByte() (byte, error)
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

const readBoolMethod = `
func readBool(r io.Reader) (bool, error) {
	var b byte
	var err error
	if br, ok := r.(ByteReader); ok {
		b, err = br.ReadByte()
	} else {
		bs := make([]byte, 1)
		_, err = io.ReadFull(r, bs)
		if err != nil {
			return false, err
		}
		b = bs[0]
	}
	return b == 1, nil
}
`

type boolField struct {
	primitiveField
}

func NewBoolField(definition interface{}) *boolField {
	return &boolField{primitiveField{
		definition:         definition,
		name:               "Bool",
		goType:             "bool",
		serializerMethod:   "writeBool",
		deserializerMethod: "readBool",
	}}
}

func (s *boolField) AddSerializer(p *generator.Package) {
	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddFunction(UTIL_FILE, "", "writeBool", writeBoolMethod)
	p.AddImport(UTIL_FILE, "io")
}

func (s *boolField) AddDeserializer(p *generator.Package) {
	p.AddStruct(UTIL_FILE, "ByteReader", byteReaderInterface)
	p.AddFunction(UTIL_FILE, "", "readBool", readBoolMethod)
	p.AddImport(UTIL_FILE, "io")
}

func (s *boolField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(bool); !ok {
		return "", fmt.Errorf("Expected bool as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}
