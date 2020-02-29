package schema

type AvroType interface {
	Name() string
	GoType() string
	IsOptional() bool

	// The name of the method which writes this field onto the wire
	SerializerMethod() string

	Children() []AvroType

	Attribute(name string) interface{}
	Definition(scope map[QualifiedName]interface{}) (interface{}, error)
	DefaultValue(lvalue string, rvalue interface{}) (string, error)

	WrapperType() string
	IsReadableBy(f AvroType, visited map[QualifiedName]interface{}) bool
}
