package types

import (
	"io"
)

type nullField struct {
	primitiveField
}

func NewNullField(definition interface{}) *nullField {
	return &nullField{primitiveField{
		definition:         definition,
		name:               "Null",
		goType:             "interface{}",
		serializerMethod:   "types.WriteNull",
		deserializerMethod: "types.ReadNull",
	}}
}

func (s *nullField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	return "", nil
}

func (s *nullField) Skip(_ io.Reader) error {
	return nil
}
