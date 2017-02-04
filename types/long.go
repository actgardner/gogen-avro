package types

import (
	"github.com/alanctgardner/gogen-avro/generator"
)

const writeLongMethod = `
func writeLong(r int64, w io.Writer) error {
	downShift := uint64(63)
	encoded := uint64((r << 1) ^ (r >> downShift))
	const maxByteSize = 10
	return encodeInt(w, maxByteSize, encoded)
}
`

const readLongMethod = `
func readLong(r io.Reader) (int64, error) {
	var v uint64
	buf := make([]byte, 1)
	for shift := uint(0); ; shift += 7 {
		if _, err := io.ReadFull(r, buf); err != nil {
			return 0, err
		}
		b := buf[0]
		v |= uint64(b&127) << shift
		if b&128 == 0 {
			break
		}
	}
	datum := (int64(v>>1) ^ -int64(v&1))
	return datum, nil
}
`

type longField struct {
	name         string
	defaultValue int64
	hasDefault   bool
}

func (s *longField) AvroName() string {
	return s.name
}

func (s *longField) GoName() string {
	return generator.ToPublicName(s.name)
}

func (s *longField) HasDefault() bool {
	return s.hasDefault
}

func (s *longField) Default() interface{} {
	return s.defaultValue
}

func (s *longField) FieldType() string {
	return "Long"
}

func (s *longField) GoType() string {
	return "int64"
}

func (s *longField) SerializerMethod() string {
	return "writeLong"
}

func (s *longField) DeserializerMethod() string {
	return "readLong"
}

func (s *longField) AddStruct(p *generator.Package) {}

func (s *longField) AddSerializer(p *generator.Package) {
	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddFunction(UTIL_FILE, "", "writeLong", writeLongMethod)
	p.AddFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.AddImport(UTIL_FILE, "io")
}

func (s *longField) AddDeserializer(p *generator.Package) {
	p.AddFunction(UTIL_FILE, "", "readLong", readLongMethod)
	p.AddImport(UTIL_FILE, "io")
}

func (s *longField) ResolveReferences(n *Namespace) error {
	return nil
}

func (s *longField) Schema(names map[QualifiedName]interface{}) interface{} {
	return "long"
}
