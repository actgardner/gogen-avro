package types

import (
	"fmt"
	"io"
)

type boolField struct {
	primitiveField
}

func NewBoolField(definition interface{}) *boolField {
	return &boolField{primitiveField{
		definition:         definition,
		name:               "Bool",
		goType:             "bool",
		serializerMethod:   "types.WriteBool",
		deserializerMethod: "types.ReadBool",
	}}
}

func (s *boolField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(bool); !ok {
		return "", fmt.Errorf("Expected bool as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}

func (s *boolField) Skip(r io.Reader) error {
	_, err := ReadBool(r)
	return err
}
