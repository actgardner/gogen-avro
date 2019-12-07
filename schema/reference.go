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

func (s *Reference) IsReadableBy(f AvroType) bool {
	if reader, ok := f.(*Reference); ok {
		return s.Def.IsReadableBy(reader.Def)
	}
	return false
}
