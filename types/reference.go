package types

import (
	"fmt"
	"github.com/alanctgardner/gogen-avro/generator"
)

/*
  A named Reference to a user-defined type (fixed, enum, record). Just a wrapper with a name around a Definition.
*/

type Reference struct {
	name         string
	typeName     QualifiedName
	def          Definition
	defaultValue interface{}
	hasDefault   bool
}

func (s *Reference) HasDefault() bool {
	return s.hasDefault
}

func (s *Reference) Default() interface{} {
	return s.defaultValue
}

func (s *Reference) Definition() Definition {
	return s.def
}

func (s *Reference) AvroName() string {
	return s.name
}

func (s *Reference) GoName() string {
	return generator.ToPublicName(s.name)
}

func (s *Reference) GoType() string {
	return s.def.GoType()
}

func (s *Reference) FieldType() string {
	return s.def.FieldType()
}

func (s *Reference) SerializerMethod() string {
	return s.def.SerializerMethod()
}

func (s *Reference) DeserializerMethod() string {
	return s.def.DeserializerMethod()
}

func (s *Reference) AddStruct(p *generator.Package) {
	s.def.AddStruct(p)
}

func (s *Reference) AddSerializer(p *generator.Package) {
	s.def.AddSerializer(p)
}

func (s *Reference) AddDeserializer(p *generator.Package) {
	s.def.AddDeserializer(p)
}

func (s *Reference) Schema(names map[QualifiedName]interface{}) interface{} {
	return s.def.Schema(names)
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
