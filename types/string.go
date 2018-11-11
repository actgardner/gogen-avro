package types

import (
	"fmt"
	"github.com/actgardner/gogen-avro/generator"
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

  // makeslice can fail depending on available memory.
  // We arbitrarily limit string size to sane default (~2.2GB).
	if len < 0 || len > math.MaxInt32 {
		return "", fmt.Errorf("string length out of range: %d", len)
	}

	if len == 0 {
		return "", nil
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
	primitiveField
}

func NewStringField(definition interface{}) *stringField {
	return &stringField{primitiveField{
		definition:         definition,
		name:               "String",
		goType:             "string",
		serializerMethod:   "writeString",
		deserializerMethod: "readString",
	}}
}

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
	p.AddImport(UTIL_FILE, "fmt")
	p.AddImport(UTIL_FILE, "math")
}

func (s *stringField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(string); !ok {
		return "", fmt.Errorf("Expected string as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %q", lvalue, rvalue), nil
}
