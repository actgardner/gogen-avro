package schema

type ArrayField struct {
	generatorMetadata

	itemType   AvroType
	definition map[string]interface{}
}

func NewArrayField(itemType AvroType, definition map[string]interface{}) *ArrayField {
	return &ArrayField{
		itemType:   itemType,
		definition: definition,
	}
}

func (s *ArrayField) ItemType() AvroType {
	return s.itemType
}

func (s *ArrayField) Attribute(name string) interface{} {
	return s.definition[name]
}

func (s *ArrayField) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	var err error
	s.definition["items"], err = s.itemType.Definition(scope)
	if err != nil {
		return nil, err
	}

	return s.definition, nil
}

func (s *ArrayField) IsReadableBy(f AvroType) bool {
	if reader, ok := f.(*ArrayField); ok {
		return s.ItemType().IsReadableBy(reader.ItemType())
	}
	return false
}

func (s *ArrayField) Children() []AvroType {
	return []AvroType{s.itemType}
}
