package schema

import (
	"fmt"
)

type DoubleField struct {
	primitiveField
}

func NewDoubleField(definition interface{}) *DoubleField {
	return &DoubleField{newPrimitiveField("Double", "float64", definition, "vm.WriteDouble")}
}

func (s *DoubleField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(float64); !ok {
		return "", fmt.Errorf("Expected number as default for field %v, got %q", lvalue, rvalue)
	}
	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}

func (s *DoubleField) WrapperType() string {
	return "types.Double"
}

func (s *DoubleField) IsReadableBy(f AvroType, visited map[QualifiedName]interface{}) bool {
	if _, ok := f.(*DoubleField); ok {
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
