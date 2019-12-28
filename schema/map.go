package schema

type MapField struct {
	itemType   AvroType
	definition map[string]interface{}
}

func NewMapField(itemType AvroType, definition map[string]interface{}) *MapField {
	return &MapField{
		itemType:   itemType,
		definition: definition,
	}
}

func (s *MapField) ItemType() AvroType {
	return s.itemType
}

func (s *MapField) Children() []AvroType {
	return []AvroType{s.itemType}
}

func (s *MapField) IsReadableBy(f AvroType) bool {
	if reader, ok := f.(*MapField); ok {
		return s.ItemType().IsReadableBy(reader.ItemType())
	}
	return false
}
