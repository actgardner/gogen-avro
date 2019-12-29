package schema

import (
	"github.com/actgardner/gogen-avro/generator"
)

type AvroType interface {
	Name() string
	GoType() string

	// The name of the method which writes this field onto the wire
	SerializerMethod() string

	// Add the imports and struct for the definition of this type to the generator.Package
	AddStruct(*generator.Package, bool) error

	Children() []AvroType

	Definition(scope map[QualifiedName]interface{}) (interface{}, error)
	DefaultValue(lvalue string, rvalue interface{}) (string, error)

	WrapperType() string
	IsReadableBy(f AvroType) bool
}
