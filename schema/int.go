package schema

type IntField struct {
	PrimitiveField
}

func NewIntField(definition interface{}) *IntField {
	return &IntField{PrimitiveField{
		definition: definition,
	}}
}

func (s *IntField) IsReadableBy(f AvroType) bool {
	if _, ok := f.(*IntField); ok {
		return true
	}
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
