package generator

type field interface {
	// The field name
	Name() string
	// The friendly type name
	FieldType() string
	// The corresponding Go type
	GoType() string
	// A method which writes this field onto the wire
	SerializerMethod() string
	// All the imports, methods and structs required for the serializer
	SerializerNs(imports, ns map[string]string)
}
