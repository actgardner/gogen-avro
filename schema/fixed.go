package schema

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/actgardner/gogen-avro/generator"
	"github.com/actgardner/gogen-avro/schema/templates"
)

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
	return generator.ToPublicName(s.name.Name)
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

func (s *FixedDefinition) SizeBytes() int {
	return s.sizeBytes
}

func (s *FixedDefinition) filename() string {
	return generator.ToSnake(s.GoType()) + ".go"
}

func (s *FixedDefinition) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.GoType())
}

func (s *FixedDefinition) structDefinition() (string, error) {
	buf := &bytes.Buffer{}
	t, err := template.New("fixed").Parse(templates.FixedTemplate)
	if err != nil {
		return "", err
	}
	err = t.Execute(buf, s)
	return buf.String(), err
}

func (s *FixedDefinition) AddStruct(p *generator.Package, _ bool) error {
	def, err := s.structDefinition()
	if err != nil {
		return err
	}

	p.AddStruct(s.filename(), s.GoType(), def)
	return nil
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

	return fmt.Sprintf("copy(%v[:], []byte(%q))", lvalue, rvalue), nil
}

func (s *FixedDefinition) IsReadableBy(d Definition) bool {
	if fixed, ok := d.(*FixedDefinition); ok {
		return fixed.sizeBytes == s.sizeBytes && fixed.name == s.name
	}
	return false
}

func (s *FixedDefinition) WrapperType() string {
	return fmt.Sprintf("%vWrapper", s.GoType())
}
