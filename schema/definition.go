package schema

/*
  The definition of a record, fixed or enum satisfies this interface.
*/

type Definition interface {
	AvroName() QualifiedName
	Aliases() []QualifiedName

	// A user-friendly name that can be built into a Go string (for unions, mostly)
	Name() string
	SimpleName() string

	GoType() string

	SerializerMethod() string

	Children() []AvroType

	// A JSON object defining this object, for writing the schema back out
	Definition(scope map[QualifiedName]interface{}) (interface{}, error)
	DefaultValue(lvalue string, rvalue interface{}) (string, error)

	IsReadableBy(f Definition) bool
	WrapperType() string
}
