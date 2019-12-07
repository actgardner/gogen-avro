package schema

import (
	"fmt"
)

type BytesField struct {
	PrimitiveField
}

func NewBytesField(definition interface{}) *BytesField {
	return &BytesField{PrimitiveField{
		definition: definition,
	}}
}

func (_ *BytesField) IsReadableBy(f AvroType) bool {
	if _, ok := f.(*BytesField); ok {
		return true
	}
	if _, ok := f.(*StringField); ok {
		return true
	}
	return false
}
