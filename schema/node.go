package schema

type Node interface {
	Children() []AvroType
	SetGeneratorMetadata(key, value interface{})
	GetGeneratorMetadata(key interface{}) interface{}
	HasGeneratorMetadata(key interface{}) bool
}
