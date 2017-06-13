package types

import (
	"fmt"
	"io"
)

type stringField struct {
	primitiveField
}

func NewStringField(definition interface{}) *stringField {
	return &stringField{primitiveField{
		definition:         definition,
		name:               "String",
		goType:             "string",
		serializerMethod:   "types.WriteString",
		deserializerMethod: "types.ReadString",
	}}
}

func (s *stringField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(string); !ok {
		return "", fmt.Errorf("Expected string as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %q", lvalue, rvalue), nil
}

func (s *stringField) Skip(r io.Reader) error {
	_, err := ReadString(r)
	return err
}
