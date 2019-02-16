package types

import (
	"fmt"

	"github.com/actgardner/gogen-avro/generator"
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
	name       QualifiedName
	aliases    []QualifiedName
	sizeBytes  int
	definition map[string]interface{}
}

func NewFixedDefinition(name QualifiedName, aliases []QualifiedName, sizeBytes int, definition map[string]interface{}) *FixedDefinition {
	return &FixedDefinition{
		name:       name,
		aliases:    aliases,
		sizeBytes:  sizeBytes,
		definition: definition,
	}
}

func (s *FixedDefinition) Name() string {
	return s.GoType()
}

func (s *FixedDefinition) SimpleName() string {
	return generator.ToPublicSimpleName(s.name.Name)
}

func (s *FixedDefinition) AvroName() QualifiedName {
	return s.name
}

func (s *FixedDefinition) Aliases() []QualifiedName {
	return s.aliases
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
	return fmt.Sprintf("write%v", s.GoType())
}

func (s *FixedDefinition) DeserializerMethod() string {
	return fmt.Sprintf("read%v", s.GoType())
}

func (s *FixedDefinition) AddStruct(p *generator.Package, _ bool) error {
	p.AddStruct(s.filename(), s.GoType(), s.typeDef())
	return nil
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

func (s *FixedDefinition) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	if _, ok := scope[s.name]; ok {
		return s.name.String(), nil
	}
	scope[s.name] = 1
	return s.definition, nil
}

func (s *FixedDefinition) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(string); !ok {
		return "", fmt.Errorf("Expected string as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = []byte(%q)", lvalue, rvalue), nil
}
