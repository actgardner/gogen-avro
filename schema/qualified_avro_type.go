package schema

type QualifiedAvroType interface {
	AvroType
	QualifiedName() QualifiedName
	Aliases() []QualifiedName
}
