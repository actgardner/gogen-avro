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

	Attribute(name string) interface{}
	// A JSON object defining this object, for writing the schema back out
	Definition(scope map[QualifiedName]interface{}) (interface{}, error)
	DefaultValue(lvalue string, rvalue interface{}) (string, error)

	IsReadableBy(f Definition) bool
	WrapperType() string
}

func copyDefinition(x map[string]interface{}) map[string]interface{} {
	if x == nil {
		return x
	}
	x1 := make(map[string]interface{})
	for name, val := range x {
		x1[name] = val
	}
	return x1
}
