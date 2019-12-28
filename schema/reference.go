package schema

/*
  A named Reference to a user-defined type (fixed, enum, record). Just a wrapper with a name around a Definition.
*/

type Reference struct {
	TypeName QualifiedName
	Def      Definition
}

func NewReference(typeName QualifiedName) *Reference {
	return &Reference{
		TypeName: typeName,
	}
}

func (s *Reference) Children() []AvroType {
	return s.Def.Children()
}

func (s *Reference) IsReadableBy(f AvroType) bool {
	if reader, ok := f.(*Reference); ok {
		return s.Def.IsReadableBy(reader.Def)
	}
	return false
}
