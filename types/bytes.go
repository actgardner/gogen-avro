package types

import (
	"fmt"
	"io"
)

type bytesField struct {
	primitiveField
}

func NewBytesField(definition interface{}) *bytesField {
	return &bytesField{primitiveField{
		definition:         definition,
		name:               "Bytes",
		goType:             "[]byte",
		serializerMethod:   "types.WriteBytes",
		deserializerMethod: "types.ReadBytes",
	}}
}

func (s *bytesField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(string); !ok {
		return "", fmt.Errorf("Expected string as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = []byte(%q)", lvalue, rvalue), nil
}

func (s *bytesField) Skip(r io.Reader) error {
	_, err := ReadBytes(r)
	return err
}
