package schema

import (
	"fmt"
)

type IntField struct {
	primitiveField
}

func NewIntField(definition interface{}) *IntField {
	return &IntField{newPrimitiveField("Int", "int32", definition, "vm.WriteInt")}
}

func (s *IntField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(float64); !ok {
		return "", fmt.Errorf("Expected number as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}

func (s *IntField) WrapperType() string {
	return "types.Int"
}

func (s *IntField) IsReadableBy(f AvroType, visited map[QualifiedName]interface{}) bool {
	if _, ok := f.(*IntField); ok {
		return true
	}
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
