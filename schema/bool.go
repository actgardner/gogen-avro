package schema

import (
	"fmt"
)

type BoolField struct {
	primitiveField
}

func NewBoolField(definition interface{}) *BoolField {
	return &BoolField{newPrimitiveField("Bool", "bool", definition, "vm.WriteBool")}
}

func (s *BoolField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(bool); !ok {
		return "", fmt.Errorf("Expected bool as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}

func (s *BoolField) WrapperType() string {
	return "types.Boolean"
}

func (s *BoolField) IsReadableBy(f AvroType, visited map[QualifiedName]interface{}) bool {
	if _, ok := f.(*BoolField); ok {
		return true
	}
	if s.primitiveField.IsReadableBy(f, visited) {
		return true
	}
	if union, ok := f.(*UnionField); ok {
		return isReadableByUnion(s, union, visited)
	}
	return false
}
