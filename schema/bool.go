package schema

import (
	"fmt"
)

type BoolField struct {
	PrimitiveField
}

func NewBoolField(definition interface{}) *BoolField {
	return &BoolField{PrimitiveField{
		definition: definition,
		name:       "Bool",
	}}
}

func (_ *BoolField) IsReadableBy(f AvroType) bool {
	_, ok := f.(*BoolField)
	return ok
}
