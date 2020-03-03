package schema

import (
	"fmt"
	"strings"

	"github.com/actgardner/gogen-avro/generator"
)

type EnumDefinition struct {
	qualifiedField
	symbols []string
	doc     string
}

func NewEnumDefinition(qname QualifiedName, aliases []QualifiedName, symbols []string, doc string, definition interface{}) *EnumDefinition {
	return &EnumDefinition{
		qualifiedField: newQualifiedField(qname, aliases, definition),
		symbols:        symbols,
		doc:            doc,
	}
}

func (e *EnumDefinition) Doc() string {
	return e.doc
}

func (e *EnumDefinition) Symbols() []string {
	return e.symbols
}

func (e *EnumDefinition) SymbolName(symbol string) string {
	return generator.ToPublicName(e.GoType() + strings.Title(symbol))
}

func (s *EnumDefinition) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(string); !ok {
		return "", fmt.Errorf("Expected string as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %v", lvalue, generator.ToPublicName(s.GoType()+strings.Title(rvalue.(string)))), nil
}

func (s *EnumDefinition) IsReadableBy(d AvroType, visited map[QualifiedName]interface{}) bool {
	otherEnum, ok := d.(*EnumDefinition)
	return ok && otherEnum.name == s.name
}

func (s *EnumDefinition) WrapperType() string {
	return "types.Int"
}
