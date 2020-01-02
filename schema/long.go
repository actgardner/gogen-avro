package schema

type LongField struct {
	PrimitiveField
}

func NewLongField(definition interface{}) *LongField {
	return &LongField{PrimitiveField{
		definition: definition,
	}}
}

func (s *LongField) IsReadableBy(f AvroType) bool {
	if _, ok := f.(*LongField); ok {
		return true
	}
	if _, ok := f.(*FloatField); ok {
		return true
	}
	if _, ok := f.(*DoubleField); ok {
		return true
	}
	return false
}
