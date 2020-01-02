package schema

type MapField struct {
	generatorMetadata

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

func (s *MapField) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	var err error
	s.definition["values"], err = s.itemType.Definition(scope)
	if err != nil {
		return nil, err
	}
	return s.definition, nil
}

func (s *MapField) IsReadableBy(f AvroType) bool {
	if reader, ok := f.(*MapField); ok {
		return s.ItemType().IsReadableBy(reader.ItemType())
	}
	return false
}

func (s *MapField) Children() []AvroType {
	return []AvroType{s.itemType}
}
