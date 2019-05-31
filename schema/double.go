package schema

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

type DoubleField struct {
	PrimitiveField
}

func NewDoubleField(definition interface{}) *DoubleField {
	return &DoubleField{PrimitiveField{
		definition:       definition,
		name:             "Double",
		goType:           "float64",
		serializerMethod: "writeDouble",
	}}
}

func (s *DoubleField) AddSerializer(p *generator.Package) {
	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddFunction(UTIL_FILE, "", "writeDouble", writeDoubleMethod)
	p.AddFunction(UTIL_FILE, "", "encodeFloat", encodeFloatMethod)
	p.AddImport(UTIL_FILE, "io")
	p.AddImport(UTIL_FILE, "math")
}

func (s *DoubleField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(float64); !ok {
		return "", fmt.Errorf("Expected number as default for field %v, got %q", lvalue, rvalue)
	}
	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}

func (s *DoubleField) WrapperType() string {
	return "types.Double"
}

func (s *DoubleField) IsReadableBy(f AvroType) bool {
	if _, ok := f.(*DoubleField); ok {
		return true
	}
	return false
}
