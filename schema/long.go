package schema

import (
	"fmt"
)

type LongField struct {
	primitiveField
}

func NewLongField(definition interface{}) *LongField {
	return &LongField{newPrimitiveField("Long", "int64", definition, "vm.WriteLong")}
}

func (s *LongField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(float64); !ok {
		return "", fmt.Errorf("Expected number as default for Field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}

func (s *LongField) WrapperType() string {
	return "types.Long"
}

func (s *LongField) IsReadableBy(f AvroType, visited map[QualifiedName]interface{}) bool {
	if _, ok := f.(*LongField); ok {
		return true
	}
	if _, ok := f.(*FloatField); ok {
		return true
	}
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
