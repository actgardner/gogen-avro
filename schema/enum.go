package schema

type EnumDefinition struct {
	generatorMetadata

	name       QualifiedName
	aliases    []QualifiedName
	symbols    []string
	doc        string
	definition map[string]interface{}
}

func NewEnumDefinition(name QualifiedName, aliases []QualifiedName, symbols []string, doc string, definition map[string]interface{}) *EnumDefinition {
	return &EnumDefinition{
		name:       name,
		aliases:    aliases,
		symbols:    symbols,
		doc:        doc,
		definition: definition,
	}
}

func (e *EnumDefinition) Doc() string {
	return e.doc
}

func (e *EnumDefinition) AvroName() QualifiedName {
	return e.name
}

func (e *EnumDefinition) Aliases() []QualifiedName {
	return e.aliases
}

func (e *EnumDefinition) Symbols() []string {
	return e.symbols
}

func (s *EnumDefinition) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	if _, ok := scope[s.name]; ok {
		return s.name.String(), nil
	}
	scope[s.name] = 1
	return s.definition, nil
}

func (s *EnumDefinition) IsReadableBy(d Definition) bool {
	otherEnum, ok := d.(*EnumDefinition)
	return ok && otherEnum.name == s.name
}

func (s *EnumDefinition) Children() []AvroType {
	return []AvroType{}
}
