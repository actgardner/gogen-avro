package schema

import (
	"fmt"
)

type DoubleField struct {
	PrimitiveField
}

func NewDoubleField(definition interface{}) *DoubleField {
	return &DoubleField{PrimitiveField{
		definition: definition,
	}}
}

func (_ *DoubleField) IsReadableBy(f AvroType) bool {
	if _, ok := f.(*DoubleField); ok {
		return true
	}
	return false
}
