package schema

import (
	"fmt"
	"strings"

	"github.com/actgardner/gogen-avro/generator"
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

func (e *EnumDefinition) IsOptional() bool {
	return false
}

func (e *EnumDefinition) SerializerMethod() string {
	return "write" + e.GoType()
}

func (e *EnumDefinition) FromStringMethod() string {
	return "New" + e.GoType() + "Value"
}

func (e *EnumDefinition) filename() string {
	return generator.ToSnake(e.GoType()) + ".go"
}

func (s *EnumDefinition) Attribute(name string) interface{} {
	return s.definition[name]
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

func (s *EnumDefinition) IsReadableBy(d Definition, visited map[QualifiedName]interface{}) bool {
	otherEnum, ok := d.(*EnumDefinition)
	return ok && otherEnum.name == s.name
}

func (s *EnumDefinition) WrapperType() string {
	return "types.Int"
}

func (s *EnumDefinition) Children() []AvroType {
	return []AvroType{}
}
