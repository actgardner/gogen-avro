package types

import (
	"github.com/alanctgardner/gogen-avro/generator"
)

/*
 * The interface implemented by all Avro field types.
 */
type Field interface {
	// The field name in the schema definition
	AvroName() string
	// The field name in the Go struct
	GoName() string
	// The friendly type name
	FieldType() string
	// Default value for the field
	HasDefault() bool
	Default() interface{}
	// The corresponding Go type
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
	// Get the objects that will serialize to the normalized JSON schema
	Schema(names map[QualifiedName]interface{}) interface{}
}
