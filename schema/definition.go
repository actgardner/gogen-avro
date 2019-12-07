package schema

import (
	"github.com/actgardner/gogen-avro/generator"
)

/*
  The definition of a record, fixed or enum satisfies this interface.
*/

type Definition interface {
	AvroName() QualifiedName
	Aliases() []QualifiedName

	// Resolve references to user-defined types
	ResolveReferences(*Namespace) error

	// A JSON object defining this object, for writing the schema back out
	Definition(scope map[QualifiedName]interface{}) (interface{}, error)
	IsReadableBy(f Definition) bool
	WrapperType() string
}
