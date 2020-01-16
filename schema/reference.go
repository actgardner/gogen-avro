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

func (s *Reference) Name() string {
	return s.Def.Name()
}

func (s *Reference) SimpleName() string {
	return s.Def.SimpleName()
}

func (s *Reference) GoType() string {
	return s.Def.GoType()
}

func (s *Reference) SerializerMethod() string {
	return s.Def.SerializerMethod()
}

func (s *Reference) Attribute(name string) interface{} {
	return s.Def.Attribute(name)
}

func (s *Reference) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	return s.Def.Definition(scope)
}

func (s *Reference) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	return s.Def.DefaultValue(lvalue, rvalue)
}

func (s *Reference) WrapperType() string {
	return s.Def.WrapperType()
}

func (s *Reference) IsReadableBy(f AvroType) bool {
	if reader, ok := f.(*Reference); ok {
		return s.Def.IsReadableBy(reader.Def)
	}
	return false
}

func (s *Reference) Children() []AvroType {
	// References can only point to Definitions and thus have no children
	return []AvroType{}
}
