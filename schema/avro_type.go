package schema

type AvroType interface {
	AbstractType

	// The name of the method which writes this field onto the wire
	SerializerMethod() string

	DefaultValue(lvalue string, rvalue interface{}) (string, error)

	WrapperType() string
	IsReadableBy(f AvroType, visited map[QualifiedName]interface{}) bool
}
