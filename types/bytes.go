package types

import (
	"fmt"
	"github.com/actgardner/gogen-avro/generator"
)

const writeBytesMethod = `
func writeBytes(r []byte, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	_, err = w.Write(r)
	return err
}
`

const readBytesMethod = `
func readBytes(r io.Reader) ([]byte, error) {
	size, err := readLong(r)
	if err != nil {
		return nil, err
	}
	if size == 0 {
		return []byte{}, nil
	}
	bb := make([]byte, size)
	_, err = io.ReadFull(r, bb)
	return bb, err
}
`

type BytesField struct {
	PrimitiveField
}

func NewBytesField(definition interface{}) *BytesField {
	return &BytesField{PrimitiveField{
		definition:         definition,
		name:               "Bytes",
		goType:             "[]byte",
		serializerMethod:   "writeBytes",
		deserializerMethod: "readBytes",
	}}
}

func (s *BytesField) AddSerializer(p *generator.Package) {
	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddFunction(UTIL_FILE, "", "writeBytes", writeBytesMethod)
	p.AddFunction(UTIL_FILE, "", "writeLong", writeLongMethod)
	p.AddFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.AddImport(UTIL_FILE, "io")
}

func (s *BytesField) AddDeserializer(p *generator.Package) {
	p.AddFunction(UTIL_FILE, "", "readBytes", readBytesMethod)
	p.AddFunction(UTIL_FILE, "", "readLong", readLongMethod)
	p.AddImport(UTIL_FILE, "io")
}

func (s *BytesField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(string); !ok {
		return "", fmt.Errorf("Expected string as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = []byte(%q)", lvalue, rvalue), nil
}
