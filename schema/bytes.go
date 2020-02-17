package schema

import (
	"fmt"
)

type BytesField struct {
	PrimitiveField
}

func NewBytesField(definition interface{}) *BytesField {
	return &BytesField{PrimitiveField{
		definition:       definition,
		name:             "Bytes",
		goType:           "[]byte",
		serializerMethod: "vm.WriteBytes",
	}}
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

func (s *BytesField) IsReadableBy(f AvroType, _ map[QualifiedName]interface{}) bool {
	if _, ok := f.(*BytesField); ok {
		return true
	}
	if _, ok := f.(*StringField); ok {
		return true
	}
	return false
}
