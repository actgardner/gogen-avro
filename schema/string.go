package schema

import (
	"fmt"
)

type StringField struct {
	primitiveField
}

func NewStringField(definition interface{}) *StringField {
	return &StringField{newPrimitiveField("String", "string", definition, "vm.WriteString")}
}

func (s *StringField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(string); !ok {
		return "", fmt.Errorf("Expected string as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %q", lvalue, rvalue), nil
}

func (s *StringField) WrapperType() string {
	return "types.String"
}

func (s *StringField) IsReadableBy(f AvroType, visited map[QualifiedName]interface{}) bool {
	if _, ok := f.(*BytesField); ok {
		return true
	}
	if _, ok := f.(*StringField); ok {
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
