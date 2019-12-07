package schema

import (
	"fmt"
)

type StringField struct {
	PrimitiveField
}

func NewStringField(definition interface{}) *StringField {
	return &StringField{PrimitiveField{
		definition: definition,
	}}
}

func (_ *StringField) IsReadableBy(f AvroType) bool {
	if _, ok := f.(*BytesField); ok {
		return true
	}
	if _, ok := f.(*StringField); ok {
		return true
	}
	return false
}
