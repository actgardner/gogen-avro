package types

import (
	"fmt"

	"github.com/actgardner/gogen-avro/generator"
)

/*
  A named Reference to a user-defined type (fixed, enum, record). Just a wrapper with a name around a Definition.
*/

type Reference struct {
	typeName QualifiedName
	def      Definition
}

func NewReference(typeName QualifiedName) *Reference {
	return &Reference{
		typeName: typeName,
	}
}

func (s *Reference) Name() string {
	return s.def.Name()
}

func (s *Reference) SimpleName() string {
	return s.def.SimpleName()
}

func (s *Reference) GoType() string {
	return s.def.GoType()
}

func (s *Reference) SerializerMethod() string {
	return s.def.SerializerMethod()
}

func (s *Reference) DeserializerMethod() string {
	return s.def.DeserializerMethod()
}

func (s *Reference) AddStruct(p *generator.Package, containers bool) error {
	return s.def.AddStruct(p, containers)
}

func (s *Reference) AddSerializer(p *generator.Package) {
	s.def.AddSerializer(p)
}

func (s *Reference) AddDeserializer(p *generator.Package) {
	s.def.AddDeserializer(p)
}

func (s *Reference) ResolveReferences(n *Namespace) error {
	if s.def == nil {
		var ok bool
		if s.def, ok = n.Definitions[s.typeName]; !ok {
			return fmt.Errorf("Unable to resolve definition of type %v", s.typeName)
		}
		return s.def.ResolveReferences(n)
	}
	return nil
}

func (s *Reference) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	return s.def.Definition(scope)
}

func (s *Reference) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	return s.def.DefaultValue(lvalue, rvalue)
}
