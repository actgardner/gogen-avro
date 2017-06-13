package types

import (
	"fmt"
	"io"
)

type doubleField struct {
	primitiveField
}

func NewDoubleField(definition interface{}) *doubleField {
	return &doubleField{primitiveField{
		definition:         definition,
		name:               "Double",
		goType:             "float64",
		serializerMethod:   "types.WriteDouble",
		deserializerMethod: "types.ReadDouble",
	}}
}

func (s *doubleField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(float64); !ok {
		return "", fmt.Errorf("Expected number as default for field %v, got %q", lvalue, rvalue)
	}
	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}

func (s *doubleField) Skip(r io.Reader) error {
	_, err := ReadDouble(r)
	return err
}
