package types

import (
	"github.com/alanctgardner/gogen-avro/generator"
)

type AvroType interface {
	Name() string
	GoType() string

	// The name of the method which writes this field onto the wire
	SerializerMethod() string
	// The name of the method which reads this field off the wire
	DeserializerMethod() string

	// Add the imports and struct for the definition of this type to the generator.Package
	AddStruct(*generator.Package)
	// Add the imports, methods and structs required for the serializer to the generator.Package
	AddSerializer(*generator.Package)
	// Add the imports, methods and structs required for the deserializer to the generator.Package
	AddDeserializer(*generator.Package)

	// Attempt to resolve references to named structs, enums or fixed fields
	ResolveReferences(*Namespace) error

	Definition(scope map[QualifiedName]interface{}) interface{}
	DefaultValue(lvalue string, rvalue interface{}) string
}
