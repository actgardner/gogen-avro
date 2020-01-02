package schema

type FloatField struct {
	PrimitiveField
}

func NewFloatField(definition interface{}) *FloatField {
	return &FloatField{PrimitiveField{
		definition: definition,
	}}
}

func (s *FloatField) IsReadableBy(f AvroType) bool {
	if _, ok := f.(*FloatField); ok {
		return true
	}
	if _, ok := f.(*DoubleField); ok {
		return true
	}
	return false
}
