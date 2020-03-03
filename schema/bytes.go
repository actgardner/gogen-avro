package schema

import (
	"fmt"
)

type BytesField struct {
	primitiveField
}

func NewBytesField(definition interface{}) *BytesField {
	return &BytesField{newPrimitiveField("Bytes", "[]byte", definition, "vm.WriteBytes")}
}

func (s *BytesField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(string); !ok {
		return "", fmt.Errorf("Expected string as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = []byte(%q)", lvalue, rvalue), nil
}

func (s *BytesField) WrapperType() string {
	return "types.Bytes"
}

func (s *BytesField) IsReadableBy(f AvroType, visited map[QualifiedName]interface{}) bool {
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
