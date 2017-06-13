package types

import (
	"fmt"
	"io"
)

type intField struct {
	primitiveField
}

func NewIntField(definition interface{}) *intField {
	return &intField{primitiveField{
		definition:         definition,
		name:               "Int",
		goType:             "int32",
		serializerMethod:   "types.WriteInt",
		deserializerMethod: "types.ReadInt",
	}}
}

func (s *intField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(float64); !ok {
		return "", fmt.Errorf("Expected number as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}

func (s *intField) Skip(r io.Reader) error {
	_, err := ReadInt(r)
	return err
}
