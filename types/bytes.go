package types

import (
	"fmt"
	"github.com/alanctgardner/gogen-avro/generator"
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
	bb := make([]byte, size)
	_, err = io.ReadFull(r, bb)
	return bb, err
}
`

type bytesField struct {
	definition interface{}
}

func (s *bytesField) Name() string {
	return "Bytes"
}

func (s *bytesField) GoType() string {
	return "[]byte"
}

func (s *bytesField) SerializerMethod() string {
	return "writeBytes"
}

func (s *bytesField) DeserializerMethod() string {
	return "readBytes"
}

func (s *bytesField) AddStruct(*generator.Package) {}

func (s *bytesField) AddSerializer(p *generator.Package) {
	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddFunction(UTIL_FILE, "", "writeBytes", writeBytesMethod)
	p.AddFunction(UTIL_FILE, "", "writeLong", writeLongMethod)
	p.AddFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.AddImport(UTIL_FILE, "io")
}

func (s *bytesField) AddDeserializer(p *generator.Package) {
	p.AddFunction(UTIL_FILE, "", "readBytes", readBytesMethod)
	p.AddFunction(UTIL_FILE, "", "readLong", readLongMethod)
	p.AddImport(UTIL_FILE, "io")
}

func (s *bytesField) ResolveReferences(n *Namespace) error {
	return nil
}

func (s *bytesField) Definition(_ map[QualifiedName]interface{}) interface{} {
	return s.definition
}

func (s *bytesField) DefaultValue(lvalue string, rvalue interface{}) string {
	return fmt.Sprintf("%v = []byte(%q)", lvalue, rvalue)
}
