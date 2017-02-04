package types

import (
	"github.com/alanctgardner/gogen-avro/generator"
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
	name         string
	defaultValue float64
	hasDefault   bool
}

func (s *doubleField) AvroName() string {
	return s.name
}

func (s *doubleField) GoName() string {
	return generator.ToPublicName(s.name)
}

func (s *doubleField) HasDefault() bool {
	return s.hasDefault
}

func (s *doubleField) Default() interface{} {
	return s.defaultValue
}

func (s *doubleField) FieldType() string {
	return "Double"
}

func (s *doubleField) GoType() string {
	return "float64"
}

func (s *doubleField) SerializerMethod() string {
	return "writeDouble"
}

func (s *doubleField) DeserializerMethod() string {
	return "readDouble"
}

func (s *doubleField) AddStruct(*generator.Package) {}

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

func (s *doubleField) ResolveReferences(n *Namespace) error {
	return nil
}

func (s *doubleField) Schema(names map[QualifiedName]interface{}) interface{} {
	return "double"
}
