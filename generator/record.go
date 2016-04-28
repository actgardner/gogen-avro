package generator

import (
	"fmt"
)

type recordField struct {
	name     string
	typeName string
	def      *recordDefinition
}

func (s *recordField) Name() string {
	return toPublicName(s.name)
}

func (s *recordField) GoType() string {
	return fmt.Sprintf("*%v", toPublicName(s.typeName))
}

func (s *recordField) FieldType() string {
	return s.typeName
}

func (s *recordField) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.typeName)
}

func (s *recordField) DeserializerMethod() string {
	return fmt.Sprintf("read%v", s.typeName)
}

/* If the record type is defined inline, add the definition to the Package */
func (s *recordField) AddStruct(p *Package) {
	if s.def != nil {
		s.def.AddStruct(p)
	}
}

func (s *recordField) AddSerializer(p *Package) {
	if s.def != nil {
		s.def.AddSerializer(p)
	}
}

func (s *recordField) AddDeserializer(p *Package) {
	if s.def != nil {
		s.def.AddDeserializer(p)
	}
}
