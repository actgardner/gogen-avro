package schema

import (
	"fmt"
)

type BoolField struct {
	PrimitiveField
}

func NewBoolField(definition interface{}) *BoolField {
	return &BoolField{PrimitiveField{
		definition:       definition,
		name:             "Bool",
		goType:           "bool",
		serializerMethod: "vm.WriteBool",
	}}
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

func (s *BoolField) IsReadableBy(f AvroType, _ map[QualifiedName]interface{}) bool {
	_, ok := f.(*BoolField)
	return ok
}
