package schema

import (
	"fmt"
	"github.com/actgardner/gogen-avro/generator"
)

const writeLongMethod = `
func writeLong(r int64, w io.Writer) error {
	downShift := uint64(63)
	encoded := uint64((r << 1) ^ (r >> downShift))
	const maxByteSize = 10
	return encodeInt(w, maxByteSize, encoded)
}
`

type LongField struct {
	PrimitiveField
}

func NewLongField(definition interface{}) *LongField {
	return &LongField{PrimitiveField{
		definition:       definition,
		name:             "Long",
		goType:           "int64",
		serializerMethod: "writeLong",
	}}
}

func (s *LongField) AddSerializer(p *generator.Package) {
	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddFunction(UTIL_FILE, "", "writeLong", writeLongMethod)
	p.AddFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.AddImport(UTIL_FILE, "io")
}

func (s *LongField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(float64); !ok {
		return "", fmt.Errorf("Expected number as default for Field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}

func (s *LongField) IsReadableBy(f AvroType) bool {
	if _, ok := f.(*LongField); ok {
		return true
	}
	if _, ok := f.(*FloatField); ok {
		return true
	}
	if _, ok := f.(*DoubleField); ok {
		return true
	}
	return false
}
