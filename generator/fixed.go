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

func (s *fixedField) typeDef() string {
	return fmt.Sprintf("type %v [%v]byte\n", s.GoType(), s.sizeBytes)
}

func (s *fixedField) SerializerNs(imports, aux map[string]string) {
	aux[s.SerializerMethod()] = s.serializerMethodDef()
	aux["ByteWriter"] = byteWriterInterface
	aux[s.GoType()] = s.typeDef()
}

func (s *fixedField) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.FieldType())
}
