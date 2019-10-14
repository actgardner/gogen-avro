package schema

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/actgardner/gogen-avro/generator"
	"github.com/actgardner/gogen-avro/schema/templates"
)

type EnumDefinition struct {
	name       QualifiedName
	aliases    []QualifiedName
	symbols    []string
	doc        string
	definition map[string]interface{}
}

func NewEnumDefinition(name QualifiedName, aliases []QualifiedName, symbols []string, doc string, definition map[string]interface{}) *EnumDefinition {
	return &EnumDefinition{
		name:       name,
		aliases:    aliases,
		symbols:    symbols,
		doc:        doc,
		definition: definition,
	}
}

func (e *EnumDefinition) Name() string {
	return e.GoType()
}

func (e *EnumDefinition) Doc() string {
	return e.doc
}

func (e *EnumDefinition) SimpleName() string {
	return e.name.Name
}

func (e *EnumDefinition) AvroName() QualifiedName {
	return e.name
}

func (e *EnumDefinition) Aliases() []QualifiedName {
	return e.aliases
}

func (e *EnumDefinition) Symbols() []string {
	return e.symbols
}

func (e *EnumDefinition) SymbolName(symbol string) string {
	return generator.ToPublicName(e.GoType() + strings.Title(symbol))
}

func (e *EnumDefinition) GoType() string {
	return generator.ToPublicName(e.name.Name)
}

func (e *EnumDefinition) SerializerMethod() string {
	return "write" + e.GoType()
}

func (e *EnumDefinition) filename() string {
	return generator.ToSnake(e.GoType()) + ".go"
}

func (e *EnumDefinition) structDefinition() (string, error) {
	buf := &bytes.Buffer{}
	t, err := template.New("enum").Parse(templates.EnumTemplate)
	if err != nil {
		return "", err
	}
	err = t.Execute(buf, e)
	return buf.String(), err
}

func (e *EnumDefinition) AddStruct(p *generator.Package, _ bool) error {
	def, err := e.structDefinition()
	if err != nil {
		return err
	}

	p.AddStruct(e.filename(), e.GoType(), def)
	return nil
}

func (s *EnumDefinition) ResolveReferences(n *Namespace) error {
	return nil
}

func (s *EnumDefinition) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	if _, ok := scope[s.name]; ok {
		return s.name.String(), nil
	}
	scope[s.name] = 1
	return s.definition, nil
}

func (s *EnumDefinition) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(string); !ok {
		return "", fmt.Errorf("Expected string as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %v", lvalue, generator.ToPublicName(s.GoType()+strings.Title(rvalue.(string)))), nil
}

func (s *EnumDefinition) IsReadableBy(d Definition) bool {
	otherEnum, ok := d.(*EnumDefinition)
	return ok && otherEnum.name == s.name
}

func (s *EnumDefinition) WrapperType() string {
	return "types.Int"
}
