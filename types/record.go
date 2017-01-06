package types

import (
	"fmt"
	"github.com/alanctgardner/gogen-avro/generator"
)

type recordField struct {
	name     string
	typeName string
	def      *RecordDefinition
}

func (s *recordField) Name() string {
	return generator.ToPublicName(s.name)
}

func (s *recordField) GoType() string {
	return fmt.Sprintf("*%v", generator.ToPublicName(s.typeName))
}

func (s *recordField) FieldType() string {
	return s.typeName
}

func (s *recordField) SerializerMethod() string {
	return fmt.Sprintf("write%v", generator.ToPublicName(s.typeName))
}

func (s *recordField) DeserializerMethod() string {
	return fmt.Sprintf("read%v", generator.ToPublicName(s.typeName))
}

/* If the record type is defined inline, add the definition to the generator.Package */
func (s *recordField) AddStruct(p *generator.Package) {
	if s.def != nil {
		s.def.AddStruct(p)
	}
}

func (s *recordField) AddSerializer(p *generator.Package) {
	if s.def != nil {
		s.def.AddSerializer(p)
	}
}

func (s *recordField) AddDeserializer(p *generator.Package) {
	if s.def != nil {
		s.def.AddDeserializer(p)
	}
}
