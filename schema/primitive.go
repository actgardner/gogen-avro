package schema

import (
	"github.com/actgardner/gogen-avro/generator"
)

// Common methods for all primitive types
type PrimitiveField struct {
	definition interface{}
}

func (s *PrimitiveField) ResolveReferences(n *Namespace) error {
	return nil
}

func (s *PrimitiveField) Definition(_ map[QualifiedName]interface{}) (interface{}, error) {
	return s.definition, nil
}
