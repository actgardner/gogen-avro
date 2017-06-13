package types

import (
	"fmt"
	"io"
)

type longField struct {
	primitiveField
}

func NewLongField(definition interface{}) *longField {
	return &longField{primitiveField{
		definition:         definition,
		name:               "Long",
		goType:             "int64",
		serializerMethod:   "types.WriteLong",
		deserializerMethod: "types.ReadLong",
	}}
}

func (s *longField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(float64); !ok {
		return "", fmt.Errorf("Expected number as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %v", lvalue, rvalue), nil
}

func (s *longField) Skip(r io.Reader) error {
	_, err := ReadLong(r)
	return err
}
