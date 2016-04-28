package generator

import (
	"fmt"
)

const writeFixedMethod = `
func %v(r %v, w io.Writer) error {
	_, err := w.Write(r[:])
	return err
}
`

const readFixedMethod = `
func %v(r io.Reader) (%v, error) {
	var bb %v
	_, err := io.ReadFull(r, bb[:])
	return bb, err
}
`

type fixedField struct {
	name         string
	typeName     string
	defaultValue []byte
	hasDefault   bool
	sizeBytes    int
}

func (s *fixedField) Name() string {
	return toPublicName(s.name)
}

func (s *fixedField) FieldType() string {
	return toPublicName(s.typeName)
}

func (s *fixedField) GoType() string {
	return s.FieldType()
}

func (s *fixedField) serializerMethodDef() string {
	return fmt.Sprintf(writeFixedMethod, s.SerializerMethod(), s.GoType())
}

func (s *fixedField) deserializerMethodDef() string {
	return fmt.Sprintf(readFixedMethod, s.DeserializerMethod(), s.GoType(), s.GoType())
}

func (s *fixedField) typeDef() string {
	return fmt.Sprintf("type %v [%v]byte\n", s.GoType(), s.sizeBytes)
}

func (s *fixedField) filename() string {
	return toSnake(s.GoType()) + ".go"
}

func (s *fixedField) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.FieldType())
}

func (s *fixedField) DeserializerMethod() string {
	return fmt.Sprintf("read%v", s.FieldType())
}

func (s *fixedField) AddStruct(p *Package) {
	p.addStruct(s.filename(), s.GoType(), s.typeDef())
}

func (s *fixedField) AddSerializer(p *Package) {
	p.addFunction(UTIL_FILE, "", s.SerializerMethod(), s.serializerMethodDef())
	p.addImport(UTIL_FILE, "io")
}

func (s *fixedField) AddDeserializer(p *Package) {
	p.addFunction(UTIL_FILE, "", s.DeserializerMethod(), s.deserializerMethodDef())
	p.addImport(UTIL_FILE, "io")
}
