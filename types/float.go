package types

import (
	"fmt"
	"io"
)

type floatField struct {
	primitiveField
}

func NewFloatField(definition interface{}) *floatField {
	return &floatField{primitiveField{
		definition:         definition,
		name:               "Float",
		goType:             "float32",
		serializerMethod:   "types.WriteFloat",
		deserializerMethod: "types.ReadFloat",
	}}
}

func (s *floatField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(float64); !ok {
		return "", fmt.Errorf("Expected float as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}

func (s *floatField) Skip(r io.Reader) error {
	_, err := ReadFloat(r)
	return err
}
