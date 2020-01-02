package schema

type NullField struct {
	PrimitiveField
}

func NewNullField(definition interface{}) *NullField {
	return &NullField{PrimitiveField{
		definition: definition,
	}}
}

func (s *NullField) IsReadableBy(f AvroType) bool {
	_, ok := f.(*NullField)
	return ok
}
