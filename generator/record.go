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

func (s *recordField) SerializerNs(imports, aux map[string]string) {
	// Records can just be a record name referencing a different schema file, or a whole nested record definition
	if s.def != nil {
		s.def.namespaceMap(imports, aux)
		aux[s.def.goName()] = s.def.structDefinition()
	}
}

func (s *recordField) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.typeName)
}
