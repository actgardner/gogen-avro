package generator

type field interface {
	// The field name
	Name() string
	// The friendly type name
	FieldType() string
	// The corresponding Go type
	GoType() string
	// The name of the method which writes this field onto the wire
	SerializerMethod() string
	// The name of the method which reads this field off the wire
	DeserializerMethod() string
	// Add the imports and struct for the definition of this type to the Package
	AddStruct(*Package)
	// Add the imports, methods and structs required for the serializer to the Package
	AddSerializer(*Package)
	// Add the imports, methods and structs required for the deserializer to the Package
	AddDeserializer(*Package)
}
