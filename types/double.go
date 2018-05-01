package types

import (
	"fmt"
	"github.com/actgardner/gogen-avro/generator"
)

const writeDoubleMethod = `
func writeDouble(r float64, w io.Writer) error {
	bits := uint64(math.Float64bits(r))
	const byteCount = 8
	return encodeFloat(w, byteCount, bits)
}
`

const readDoubleMethod = `
func readDouble(r io.Reader) (float64, error) {
	buf := make([]byte, 8)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return 0, err
	}
	bits := binary.LittleEndian.Uint64(buf)
	val := math.Float64frombits(bits)
	return val, nil
}
`

type doubleField struct {
	primitiveField
}

func NewDoubleField(definition interface{}) *doubleField {
	return &doubleField{primitiveField{
		definition:         definition,
		name:               "Double",
		goType:             "float64",
		serializerMethod:   "writeDouble",
		deserializerMethod: "readDouble",
	}}
}

func (s *doubleField) AddSerializer(p *generator.Package) {
	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddFunction(UTIL_FILE, "", "writeDouble", writeDoubleMethod)
	p.AddFunction(UTIL_FILE, "", "encodeFloat", encodeFloatMethod)
	p.AddImport(UTIL_FILE, "io")
	p.AddImport(UTIL_FILE, "math")
}

func (s *doubleField) AddDeserializer(p *generator.Package) {
	p.AddFunction(UTIL_FILE, "", "readDouble", readDoubleMethod)
	p.AddImport(UTIL_FILE, "io")
	p.AddImport(UTIL_FILE, "math")
	p.AddImport(UTIL_FILE, "encoding/binary")
}

func (s *doubleField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(float64); !ok {
		return "", fmt.Errorf("Expected number as default for field %v, got %q", lvalue, rvalue)
	}
	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}
