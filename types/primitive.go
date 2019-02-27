package types

import (
	"github.com/actgardner/gogen-avro/generator"
)

// Common methods for all primitive types
type primitiveField struct {
	definition         interface{}
	name               string
	goType             string
	serializerMethod   string
	deserializerMethod string
}

func (s *primitiveField) Name() string {
	return s.name
}

func (s *primitiveField) SimpleName() string {
	return s.name
}

func (s *primitiveField) GoType() string {
	return s.goType
}

func (s *primitiveField) SerializerMethod() string {
	return s.serializerMethod
}

func (s *primitiveField) DeserializerMethod() string {
	return s.deserializerMethod
}

func (s *primitiveField) AddStruct(p *generator.Package, _ bool) error {
	return nil
}

func (s *primitiveField) ResolveReferences(n *Namespace) error {
	return nil
}

func (s *primitiveField) Definition(_ map[QualifiedName]interface{}) (interface{}, error) {
	return s.definition, nil
}
