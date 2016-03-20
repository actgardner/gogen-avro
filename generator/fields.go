package generator

type field interface {
	// The field name
	Name() string
	// The friendly type name
	FieldType() string
	// The corresponding Go type
	GoType() string
	// Auxillary types (enums, structs) and imports to be created
	AuxStructs(types map[string]string, imports map[string]string)
	// A method which writes this field onto the wire
	SerializerMethod() string
}
