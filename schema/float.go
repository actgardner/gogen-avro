package schema

import (
	"fmt"
)

type FloatField struct {
	primitiveField
}

func NewFloatField(definition interface{}) *FloatField {
	return &FloatField{newPrimitiveField("Float", "float32", definition, "vm.WriteFloat")}
}

func (s *FloatField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(float64); !ok {
		return "", fmt.Errorf("Expected float as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}

func (s *FloatField) WrapperType() string {
	return "types.Float"
}

func (s *FloatField) IsReadableBy(f AvroType, visited map[QualifiedName]interface{}) bool {
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
