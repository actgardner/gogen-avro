package schema

import (
	"fmt"
)

type FixedDefinition struct {
	qualifiedField
	sizeBytes int
}

func NewFixedDefinition(qname QualifiedName, aliases []QualifiedName, sizeBytes int, definition interface{}) *FixedDefinition {
	return &FixedDefinition{
		qualifiedField: newQualifiedField(qname, aliases, definition),
		sizeBytes:      sizeBytes,
	}
}

func (s *FixedDefinition) SizeBytes() int {
	return s.sizeBytes
}

func (s *FixedDefinition) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(string); !ok {
		return "", fmt.Errorf("Expected string as default for field %v, got %q", lvalue, rvalue)
	}
	return fmt.Sprintf("copy(%v[:], []byte(%q))", lvalue, rvalue), nil
}

func (s *FixedDefinition) IsReadableBy(d AvroType, visited map[QualifiedName]interface{}) bool {
	if fixed, ok := d.(*FixedDefinition); ok {
		return fixed.sizeBytes == s.sizeBytes && fixed.name == s.name
	}
	return false
}
