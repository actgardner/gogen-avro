package schema

type FixedDefinition struct {
	generatorMetadata

	name       QualifiedName
	aliases    []QualifiedName
	sizeBytes  int
	definition map[string]interface{}
}

func NewFixedDefinition(name QualifiedName, aliases []QualifiedName, sizeBytes int, definition map[string]interface{}) *FixedDefinition {
	return &FixedDefinition{
		name:       name,
		aliases:    aliases,
		sizeBytes:  sizeBytes,
		definition: definition,
	}
}

func (s *FixedDefinition) AvroName() QualifiedName {
	return s.name
}

func (s *FixedDefinition) Aliases() []QualifiedName {
	return s.aliases
}

func (s *FixedDefinition) SizeBytes() int {
	return s.sizeBytes
}

func (s *FixedDefinition) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	if _, ok := scope[s.name]; ok {
		return s.name.String(), nil
	}
	scope[s.name] = 1
	return s.definition, nil
}

func (s *FixedDefinition) IsReadableBy(d Definition) bool {
	if fixed, ok := d.(*FixedDefinition); ok {
		return fixed.sizeBytes == s.sizeBytes && fixed.name == s.name
	}
	return false
}

func (s *FixedDefinition) Children() []AvroType {
	return []AvroType{}
}
