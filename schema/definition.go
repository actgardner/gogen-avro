package schema

/*
  The definition of a record, fixed or enum satisfies this interface.
*/

type Definition interface {
	Node

	AvroName() QualifiedName
	Aliases() []QualifiedName

	// A JSON object defining this object, for writing the schema back out
	Definition(scope map[QualifiedName]interface{}) (interface{}, error)

	IsReadableBy(f Definition) bool
}
