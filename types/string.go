package types

import (
	"github.com/alanctgardner/gogen-avro/generator"
)

const stringWriterInterface = `
type StringWriter interface {
	WriteString(string) (int, error)
}
`

const writeStringMethod = `
func writeString(r string, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	if sw, ok := w.(StringWriter); ok {
		_, err = sw.WriteString(r)
	} else {
		_, err = w.Write([]byte(r))
	}
	return err
}
`

const readStringMethod = `
func readString(r io.Reader) (string, error) {
	len, err := readLong(r)
	if err != nil {
		return "", err
	}
	bb := make([]byte, len)
	_, err = io.ReadFull(r, bb)
	if err != nil {
		return "", err
	}
	return string(bb), nil
}
`

type stringField struct {
	name         string
	defaultValue string
	hasDefault   bool
}

func (s *stringField) HasDefault() bool {
	return s.hasDefault
}

func (s *stringField) Default() interface{} {
	return s.defaultValue
}

func (s *stringField) AvroName() string {
	return s.name
}

func (s *stringField) GoName() string {
	return generator.ToPublicName(s.name)
}

func (s *stringField) FieldType() string {
	return "String"
}

func (s *stringField) GoType() string {
	return "string"
}

func (s *stringField) SerializerMethod() string {
	return "writeString"
}

func (s *stringField) DeserializerMethod() string {
	return "readString"
}

func (s *stringField) AddStruct(*generator.Package) {}

func (s *stringField) AddSerializer(p *generator.Package) {
	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddStruct(UTIL_FILE, "StringWriter", stringWriterInterface)
	p.AddFunction(UTIL_FILE, "", "writeLong", writeLongMethod)
	p.AddFunction(UTIL_FILE, "", "writeString", writeStringMethod)
	p.AddFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.AddImport(UTIL_FILE, "io")
}

func (s *stringField) AddDeserializer(p *generator.Package) {
	p.AddFunction(UTIL_FILE, "", "readLong", readLongMethod)
	p.AddFunction(UTIL_FILE, "", "readString", readStringMethod)
	p.AddImport(UTIL_FILE, "io")
}

func (s *stringField) ResolveReferences(n *Namespace) error {
	return nil
}

func (s *stringField) Schema(names map[QualifiedName]interface{}) interface{} {
	return "string"
}
