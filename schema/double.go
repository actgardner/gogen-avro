package schema

import (
	"fmt"
)

type DoubleField struct {
	PrimitiveField
}

func NewDoubleField(definition interface{}) *DoubleField {
	return &DoubleField{PrimitiveField{
		definition:       definition,
		name:             "Double",
		goType:           "float64",
		serializerMethod: "vm.WriteDouble",
	}}
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

func (s *DoubleField) IsReadableBy(f AvroType, _ map[QualifiedName]interface{}) bool {
	if _, ok := f.(*DoubleField); ok {
		return true
	}
	return false
}
