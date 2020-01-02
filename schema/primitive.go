package schema

// Common methods for all primitive types
type PrimitiveField struct {
	generatorMetadata

	definition interface{}
}

func (s *PrimitiveField) Attribute(name string) interface{} {
	definition, _ := s.definition.(map[string]interface{})
	return definition[name]
}

func (s *PrimitiveField) Definition(_ map[QualifiedName]interface{}) (interface{}, error) {
	return s.definition, nil
}

func (s *PrimitiveField) Children() []AvroType {
	return []AvroType{}
}
