package schema

import (
	"github.com/actgardner/gogen-avro/generator"
)

type AvroType interface {
	// Attempt to resolve references to named structs, enums or fixed fields
	ResolveReferences(*Namespace) error

	// The definition of the type, with all nested type references expanded
	Definition(scope map[QualifiedName]interface{}) (interface{}, error)

	// True if this field can be assigned a value of type f (ex. long can be assigned an int, but not vice-versa)
	IsReadableBy(f AvroType) bool
}
