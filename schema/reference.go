package schema

import (
	"fmt"

	"github.com/actgardner/gogen-avro/generator"
)

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

func (s *Reference) AddStruct(p *generator.Package, containers bool) error {
	return s.Def.AddStruct(p, containers)
}

func (s *Reference) AddSerializer(p *generator.Package) {
	s.Def.AddSerializer(p)
}

func (s *Reference) ResolveReferences(n *Namespace) error {
	if s.Def == nil {
		var ok bool
		if s.Def, ok = n.Definitions[s.TypeName]; !ok {
			return fmt.Errorf("Unable to resolve definition of type %v", s.TypeName)
		}
		return s.Def.ResolveReferences(n)
	}
	return nil
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
