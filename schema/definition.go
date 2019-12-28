package schema

/*
  The definition of a record, fixed or enum satisfies this interface.
*/

type Definition interface {
	AvroName() QualifiedName
	Aliases() []QualifiedName
	IsReadableBy(f Definition) bool
	Children() []AvroType
}
