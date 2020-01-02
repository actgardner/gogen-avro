package schema

type AvroType interface {
	Node

	Attribute(name string) interface{}
	Definition(scope map[QualifiedName]interface{}) (interface{}, error)
	IsReadableBy(f AvroType) bool
}
