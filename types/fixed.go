package types

import (
	"fmt"
	"github.com/alanctgardner/gogen-avro/generator"
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

type FixedDefinition struct {
	name      QualifiedName
	aliases   []QualifiedName
	sizeBytes int
	metadata  map[string]interface{}
}

func (s *FixedDefinition) AvroName() QualifiedName {
	return s.name
}

func (s *FixedDefinition) Aliases() []QualifiedName {
	return s.aliases
}

func (s *FixedDefinition) FieldType() string {
	return s.GoType()
}

func (s *FixedDefinition) GoType() string {
	return generator.ToPublicName(s.name.Name)
}

func (s *FixedDefinition) serializerMethodDef() string {
	return fmt.Sprintf(writeFixedMethod, s.SerializerMethod(), s.GoType())
}

func (s *FixedDefinition) deserializerMethodDef() string {
	return fmt.Sprintf(readFixedMethod, s.DeserializerMethod(), s.GoType(), s.GoType())
}

func (s *FixedDefinition) typeDef() string {
	return fmt.Sprintf("type %v [%v]byte\n", s.GoType(), s.sizeBytes)
}

func (s *FixedDefinition) filename() string {
	return generator.ToSnake(s.GoType()) + ".go"
}

func (s *FixedDefinition) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.FieldType())
}

func (s *FixedDefinition) DeserializerMethod() string {
	return fmt.Sprintf("read%v", s.FieldType())
}

func (s *FixedDefinition) AddStruct(p *generator.Package) {
	p.AddStruct(s.filename(), s.GoType(), s.typeDef())
}

func (s *FixedDefinition) AddSerializer(p *generator.Package) {
	p.AddFunction(UTIL_FILE, "", s.SerializerMethod(), s.serializerMethodDef())
	p.AddImport(UTIL_FILE, "io")
}

func (s *FixedDefinition) AddDeserializer(p *generator.Package) {
	p.AddFunction(UTIL_FILE, "", s.DeserializerMethod(), s.deserializerMethodDef())
	p.AddImport(UTIL_FILE, "io")
}

func (s *FixedDefinition) ResolveReferences(n *Namespace) error {
	return nil
}

func (s *FixedDefinition) Schema(names map[QualifiedName]interface{}) interface{} {
	name := s.name.String()
	if _, ok := names[s.name]; ok {
		return name
	}
	names[s.name] = 1
	return mergeMaps(map[string]interface{}{
		"type": "fixed",
		"name": name,
		"size": s.sizeBytes,
	}, s.metadata)
}
