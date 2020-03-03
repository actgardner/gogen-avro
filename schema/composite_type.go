package schema

// CompositeType must be implemented by any type containing children (array, record and union)
type CompositeType interface {
	Name() string
	Children() []AvroType
}
