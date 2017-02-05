package types

type PrimitiveType struct {
	AvroType                     string
	GoType                       string
	SerializerMethod             string
	DeserializerMethod           string
	SerializerMethodDefinition   string
	DeserializerMethodDefinition string
}
