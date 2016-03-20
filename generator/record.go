package generator

import (
	"fmt"
)

type recordField struct {
	name     string
	typeName string
}

func (s *recordField) Name() string {
	return toPublicName(s.name)
}

func (s *recordField) GoType() string {
	return fmt.Sprintf("*%v", toPublicName(s.typeName))
}

func (s *recordField) FieldType() string {
	return s.typeName + "Record"
}

func (s *recordField) AuxStructs(aux map[string]string, _ map[string]string) {
}

func (s *recordField) SerializerMethod() string {
	return fmt.Sprintf("write%v", s.typeName)
}
