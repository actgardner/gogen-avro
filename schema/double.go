package schema

import (
	"fmt"
	"reflect"

	"github.com/actgardner/gogen-avro/generator"
)

// FieldName represents a Avro field type name
type FieldName string

// Available field names.
// TODO: make use of `FieldName` type
const (
	NameDoubleField FieldName = "Double"
)

// Method represents a method that could be passed on to a schema generator
type Method struct {
	Name      FieldName // Name of the method
	Generator string    // Content to be included during generation
}

// WriteDouble represents a write double method that includes the content to be included during generation
// TODO: implement method(s)
var WriteDouble = Method{
	Name: NameDoubleField,
	Generator: `
func writeDouble(r float64, w io.Writer) error {
	bits := uint64(math.Float64bits(r))
	const byteCount = 8
	return encodeFloat(w, byteCount, bits)
}
`,
}

const writeDoubleMethod = `
func writeDouble(r float64, w io.Writer) error {
	bits := uint64(math.Float64bits(r))
	const byteCount = 8
	return encodeFloat(w, byteCount, bits)
}
`

// DoubleField represents a float(64) Avro field.
type DoubleField struct {
	PrimitiveField
}

// NewDoubleField constructs a new `DoubleField` for the given definition
func NewDoubleField(definition interface{}) *DoubleField {
	return &DoubleField{PrimitiveField{
		definition:       definition,
		name:             "Double",
		goType:           reflect.Float64.String(),
		serializerMethod: "writeDouble",
	}}
}

// AddSerializer includes the required methods, structs and imports for a `DoubleField` to the given generator package.
func (s *DoubleField) AddSerializer(p *generator.Package) {
	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddFunction(UTIL_FILE, "", "writeDouble", writeDoubleMethod)
	p.AddFunction(UTIL_FILE, "", "encodeFloat", encodeFloatMethod)
	p.AddImport(UTIL_FILE, "io")
	p.AddImport(UTIL_FILE, "math")
}

// DefaultValue is a generator method that returns a default value constructor.
// It expects a int64/float64 as default value for the given field and returns a error if given otherwise.
func (s *DoubleField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(float64); !ok {
		return "", fmt.Errorf("Expected number as default for field %v, got %q", lvalue, rvalue)
	}
	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}

// WrapperType returns the Avro type representation
func (s *DoubleField) WrapperType() string {
	return "types.Double"
}

// IsReadableBy preforms a check if the given `AvroType` is readable by `DoubleField`.
// The method will return false if the given type is not a `DoubleField` pointer.
func (s *DoubleField) IsReadableBy(f AvroType) bool {
	if _, ok := f.(*DoubleField); ok {
		return true
	}
	return false
}
