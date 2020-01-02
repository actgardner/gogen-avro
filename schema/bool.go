package schema

type BoolField struct {
	PrimitiveField
}

func NewBoolField(definition interface{}) *BoolField {
	return &BoolField{PrimitiveField{
		definition: definition,
	}}
}

func (s *BoolField) IsReadableBy(f AvroType) bool {
	_, ok := f.(*BoolField)
	return ok
}
