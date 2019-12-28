package schema

type AvroType interface {
	// Children returns a slice with all the AvroTypes referenced by this AvroType
	Children() []AvroType

	// True if this field can be assigned a value of type f (ex. long can be assigned an int, but not vice-versa)
	IsReadableBy(f AvroType) bool
}
